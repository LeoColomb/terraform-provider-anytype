// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIconMarshalEmoji(t *testing.T) {
	got, err := json.Marshal(Icon{Format: IconFormatEmoji, Emoji: "📄", File: "ignored"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"format":"emoji","emoji":"📄"}`
	if string(got) != want {
		t.Errorf("emoji marshal: got %s want %s", got, want)
	}
}

func TestIconMarshalFile(t *testing.T) {
	got, err := json.Marshal(Icon{Format: IconFormatFile, File: "bafy", Emoji: "ignored"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"format":"file","file":"bafy"}`
	if string(got) != want {
		t.Errorf("file marshal: got %s want %s", got, want)
	}
}

func TestIconMarshalNamed(t *testing.T) {
	got, err := json.Marshal(Icon{Format: IconFormatIcon, Name: "book", Color: "yellow", File: "ignored"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"format":"icon","name":"book","color":"yellow"}`
	if string(got) != want {
		t.Errorf("icon marshal: got %s want %s", got, want)
	}
}

func TestIconMarshalEmpty(t *testing.T) {
	got, err := json.Marshal(Icon{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(got) != "null" {
		t.Errorf("empty icon: got %s want null", got)
	}
}

func TestIconMarshalUnknownFormatErrors(t *testing.T) {
	_, err := json.Marshal(Icon{Format: "bogus"})
	if err == nil || !strings.Contains(err.Error(), "unknown icon format") {
		t.Errorf("expected unknown format error, got %v", err)
	}
}

func TestIconUnmarshalEmoji(t *testing.T) {
	var i Icon
	if err := json.Unmarshal([]byte(`{"format":"emoji","emoji":"📄"}`), &i); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if i.Format != IconFormatEmoji || i.Emoji != "📄" {
		t.Errorf("unexpected: %+v", i)
	}
}

func TestIconUnmarshalNamed(t *testing.T) {
	var i Icon
	if err := json.Unmarshal([]byte(`{"format":"icon","name":"book","color":"yellow"}`), &i); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if i.Format != IconFormatIcon || i.Name != "book" || i.Color != "yellow" {
		t.Errorf("unexpected: %+v", i)
	}
}

func TestIconUnmarshalNullAndEmpty(t *testing.T) {
	var i Icon
	if err := json.Unmarshal([]byte(`null`), &i); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if (i != Icon{}) {
		t.Errorf("null should yield zero Icon, got %+v", i)
	}
}

func TestIconUnmarshalInfersFormat(t *testing.T) {
	var i Icon
	if err := json.Unmarshal([]byte(`{"emoji":"📄"}`), &i); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if i.Format != IconFormatEmoji {
		t.Errorf("expected inferred emoji format, got %q", i.Format)
	}
}
