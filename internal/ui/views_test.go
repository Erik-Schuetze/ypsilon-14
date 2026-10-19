package ui

import (
	"strings"
	"testing"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestWrapWords(t *testing.T) {
	lines := wrapWords(titleText, 90)
	if len(lines) != 2 {
		t.Fatalf("wrapWords(title, 90) returned %d lines, want 2", len(lines))
	}
	for _, line := range lines {
		if runeLen(line) > 90 {
			t.Errorf("wrapped title line %q is %d runes, want <= 90", line, runeLen(line))
		}
	}
	joined := strings.Join(strings.Fields(strings.Join(lines, " ")), " ")
	want := strings.Join(strings.Fields(titleText), " ")
	if joined != want {
		t.Errorf("wrapped title changed its words: %q != %q", joined, want)
	}

	if got := wrapWords("", 10); len(got) != 1 || got[0] != "" {
		t.Errorf("wrapWords(\"\", 10) = %#v, want [\"\"]", got)
	}

	got := wrapWords("hello world", 6)
	if len(got) != 2 {
		t.Fatalf("wrapWords(\"hello world\", 6) returned %d lines, want 2", len(got))
	}
	for _, line := range got {
		if runeLen(line) > 6 {
			t.Errorf("wrapped line %q is %d runes, want <= 6", line, runeLen(line))
		}
	}
}

func TestWrappedLineCount(t *testing.T) {
	cases := []struct {
		text  string
		width int
		want  int
	}{
		{"", 10, 1},
		{"abc", 3, 2},
		{"abcd", 3, 2},
		{"a\nb", 5, 2},
		{"abc", 0, 1},
	}
	for _, c := range cases {
		if got := wrappedLineCount(c.text, c.width); got != c.want {
			t.Errorf("wrappedLineCount(%q, %d) = %d, want %d", c.text, c.width, got, c.want)
		}
	}
}

func TestMessageRect(t *testing.T) {
	cases := []struct {
		w, h  int
		title string
		msg   string
	}{
		{100, 29, "CARD READER v1.2", "Enabled admin access"},
		{100, 29, "SELF DESTRUCTION", "Initiated Self Destruction. Evacuate immediately!"},
		{100, 29, "ERROR", "Admin access required!"},
		{100, 29, "LONG", strings.Repeat("x", 200)},
		{100, 29, strings.Repeat("T", 200), "ok"},
		{100, 29, strings.Repeat("T", 60), "ok"},
		{160, 45, "CARD READER v1.2", "Enabled admin access"},
	}

	for _, c := range cases {
		r, innerW := messageRect(c.w, c.h, c.title, c.msg)
		if !r.Valid(c.w, c.h) {
			t.Errorf("%dx%d %q: rect %+v is not valid", c.w, c.h, c.title, r)
		}
		if r.InnerW() != innerW {
			t.Errorf("%dx%d %q: InnerW %d != returned %d", c.w, c.h, c.title, r.InnerW(), innerW)
		}
		if innerW < 1 || r.InnerH() < 1 {
			t.Errorf("%dx%d %q: no drawable inner area: %+v", c.w, c.h, c.title, r)
		}
		if need := wrappedLineCount(" "+c.msg, innerW); need > r.InnerH() {
			t.Errorf("%dx%d %q: message needs %d rows, has %d", c.w, c.h, c.title, need, r.InnerH())
		}
		// drawTitle renders at most innerW-2 runes of " "+title+" ".
		if need := runeLen(" "+c.title+" ") + 2; need <= c.w-6 && innerW < need {
			t.Errorf("%dx%d %q: title needs %d inner columns, has %d", c.w, c.h, c.title, need, innerW)
		}
		if diff := abs((r.X0+r.X1)/2 - c.w/2); diff > 1 {
			t.Errorf("%dx%d %q: horizontal center off by %d", c.w, c.h, c.title, diff)
		}
		if diff := abs((r.Y0+r.Y1)/2 - c.h/2); diff > 1 {
			t.Errorf("%dx%d %q: vertical center off by %d", c.w, c.h, c.title, diff)
		}
		if r.X0 < 2 || r.X1 > c.w-3 {
			t.Errorf("%dx%d %q: horizontal margins too small: %+v", c.w, c.h, c.title, r)
		}
	}
}

func TestMessageRectRecenters(t *testing.T) {
	r1, _ := messageRect(120, 40, "CARD READER v1.2", "Enabled admin access")
	r2, _ := messageRect(160, 45, "CARD READER v1.2", "Enabled admin access")

	if r1.X0 == r2.X0 || r1.Y0 == r2.Y0 {
		t.Errorf("popup did not move on resize: %+v vs %+v", r1, r2)
	}
	if diff := abs((r1.X0+r1.X1)/2 - 60); diff > 1 {
		t.Errorf("120-wide popup center off by %d", diff)
	}
	if diff := abs((r2.X0+r2.X1)/2 - 80); diff > 1 {
		t.Errorf("160-wide popup center off by %d", diff)
	}
}

func TestSizeWarningText(t *testing.T) {
	text := sizeWarningText(80, 24)
	for _, want := range []string{
		"TERMINAL TOO SMALL",
		"Required: 100 x 29",
		"Current:  80 x 24",
		"CTRL-C to quit.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("sizeWarningText(80,24) missing %q:\n%s", want, text)
		}
	}

	if got := sizeWarningText(2, 2); got != "" {
		t.Errorf("sizeWarningText(2,2) = %q, want empty", got)
	}
}
