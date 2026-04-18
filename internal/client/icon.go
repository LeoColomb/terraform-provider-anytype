// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Icon mirrors the polymorphic Icon schema from the Anytype OpenAPI. The API
// expresses Icon as a `oneOf` (EmojiIcon | FileIcon | NamedIcon) discriminated
// by the `format` field (one of "emoji", "file", "icon"). The Terraform Plugin
// Framework code generator does not yet support `oneOf`, so we model all three
// variants in a single struct and rely on custom JSON marshalling to keep the
// wire format faithful to the spec.
//
// Fields not relevant to the active Format are omitted on the wire.
type Icon struct {
	// Format is the discriminator; one of IconFormat* below.
	Format string `json:"format,omitempty"`
	// Emoji is set when Format == IconFormatEmoji.
	Emoji string `json:"emoji,omitempty"`
	// File is set when Format == IconFormatFile (a CID reference).
	File string `json:"file,omitempty"`
	// Name is set when Format == IconFormatIcon (one of the IconName enum).
	Name string `json:"name,omitempty"`
	// Color is set when Format == IconFormatIcon (one of the Color enum).
	Color string `json:"color,omitempty"`
}

// Icon format discriminator values, mirroring the IconFormat enum.
const (
	IconFormatEmoji = "emoji"
	IconFormatFile  = "file"
	IconFormatIcon  = "icon"
)

// MarshalJSON emits only the variant fields appropriate for Format. An empty
// Format with no variant fields set marshals as `null` so the icon can be
// omitted from request bodies through `omitempty` on a *Icon pointer.
func (i Icon) MarshalJSON() ([]byte, error) {
	switch i.Format {
	case IconFormatEmoji:
		return json.Marshal(struct {
			Format string `json:"format"`
			Emoji  string `json:"emoji,omitempty"`
		}{i.Format, i.Emoji})
	case IconFormatFile:
		return json.Marshal(struct {
			Format string `json:"format"`
			File   string `json:"file,omitempty"`
		}{i.Format, i.File})
	case IconFormatIcon:
		return json.Marshal(struct {
			Format string `json:"format"`
			Name   string `json:"name,omitempty"`
			Color  string `json:"color,omitempty"`
		}{i.Format, i.Name, i.Color})
	case "":
		// No discriminator — let the caller decide whether to send anything.
		// We still emit any populated fields so a hand-built Icon is not
		// silently dropped.
		if i.Emoji == "" && i.File == "" && i.Name == "" && i.Color == "" {
			return []byte("null"), nil
		}
		type rawIcon Icon
		return json.Marshal(rawIcon(i))
	default:
		return nil, fmt.Errorf("anytype: unknown icon format %q", i.Format)
	}
}

// UnmarshalJSON decodes any of the three Icon variants. The `format`
// discriminator is required by the API, but we tolerate its absence by
// inferring it from whichever variant field is populated.
func (i *Icon) UnmarshalJSON(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 || bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*i = Icon{}
		return nil
	}
	var raw struct {
		Format string `json:"format"`
		Emoji  string `json:"emoji"`
		File   string `json:"file"`
		Name   string `json:"name"`
		Color  string `json:"color"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*i = Icon{
		Format: raw.Format,
		Emoji:  raw.Emoji,
		File:   raw.File,
		Name:   raw.Name,
		Color:  raw.Color,
	}
	if i.Format == "" {
		switch {
		case i.Emoji != "":
			i.Format = IconFormatEmoji
		case i.File != "":
			i.Format = IconFormatFile
		case i.Name != "" || i.Color != "":
			i.Format = IconFormatIcon
		}
	}
	return nil
}
