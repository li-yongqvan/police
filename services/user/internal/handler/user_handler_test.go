package handler

import (
	"strings"
	"testing"
	"time"
)

func TestAvatarFilenameUsesNanosecondPrecision(t *testing.T) {
	base := time.Date(2026, 7, 20, 14, 0, 0, 0, time.UTC)

	first := avatarFilename(42, "avatar.png", base)
	second := avatarFilename(42, "avatar.png", base.Add(time.Nanosecond))

	if first == second {
		t.Fatalf("avatar filenames should be unique for consecutive uploads: %q", first)
	}
	if !strings.HasPrefix(first, "42_") || !strings.HasSuffix(first, ".png") {
		t.Fatalf("avatar filename should preserve user prefix and extension: %q", first)
	}
}
