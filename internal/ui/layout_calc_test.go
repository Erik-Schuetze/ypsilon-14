package ui

import "testing"

func rectsOf(l Layout) []struct {
	name string
	r    Rect
} {
	return []struct {
		name string
		r    Rect
	}{
		{"title", l.Title},
		{"blink", l.Blink},
		{"menu", l.Menu},
		{"key", l.Key},
		{"log", l.Log},
		{"main", l.Main},
		{"legend", l.Legend},
		{"operation", l.Operation},
	}
}

func checkLayoutInvariants(t *testing.T, w, h int, l Layout) {
	t.Helper()

	for _, entry := range rectsOf(l) {
		if !entry.r.Valid(w, h) {
			t.Errorf("%dx%d: %s rect %+v is not valid", w, h, entry.name, entry.r)
		}
	}

	if l.Main.InnerW() < MainTextWidth {
		t.Errorf("%dx%d: main inner width %d < %d", w, h, l.Main.InnerW(), MainTextWidth)
	}
	if l.Main.InnerH() < FloorplanHeight {
		t.Errorf("%dx%d: main inner height %d < %d", w, h, l.Main.InnerH(), FloorplanHeight)
	}
	if l.Menu.InnerW() < MenuTextWidth {
		t.Errorf("%dx%d: menu inner width %d < %d", w, h, l.Menu.InnerW(), MenuTextWidth)
	}
	if l.Menu.InnerH() < 5 {
		t.Errorf("%dx%d: menu inner height %d < 5", w, h, l.Menu.InnerH())
	}
	if l.Legend.InnerW() < LegendTextWidth {
		t.Errorf("%dx%d: legend inner width %d < %d", w, h, l.Legend.InnerW(), LegendTextWidth)
	}
	if l.Legend.InnerH() < LegendHeight {
		t.Errorf("%dx%d: legend inner height %d < %d", w, h, l.Legend.InnerH(), LegendHeight)
	}
	if l.Operation.InnerW() < OperationTextWidth {
		t.Errorf("%dx%d: operation inner width %d < %d", w, h, l.Operation.InnerW(), OperationTextWidth)
	}
	if l.Operation.InnerH() < OperationHeight {
		t.Errorf("%dx%d: operation inner height %d < %d", w, h, l.Operation.InnerH(), OperationHeight)
	}
	if l.Log.InnerW() < 1 || l.Log.InnerH() < 1 {
		t.Errorf("%dx%d: log has no drawable inner area: %+v", w, h, l.Log)
	}
	if l.Blink.InnerW() < 3 {
		t.Errorf("%dx%d: blink inner width %d < 3", w, h, l.Blink.InnerW())
	}

	// The overlay panels must leave enough of the main pane visible for the
	// fixed floorplan / operations list.
	legendVisible := l.Legend.X0 - (l.Main.X0 + 1)
	if legendVisible < FloorplanWidth {
		t.Errorf("%dx%d: main visible under legend = %d < %d", w, h, legendVisible, FloorplanWidth)
	}
	operationVisible := l.Operation.X0 - (l.Main.X0 + 1)
	if operationVisible < OperationsListWidth {
		t.Errorf("%dx%d: main visible under operations = %d < %d", w, h, operationVisible, OperationsListWidth)
	}

	// Title and navigation hint must fit in the allocated rows.
	titleLines := wrapWords(titleText, l.Title.InnerW())
	if len(titleLines) > l.Title.InnerH() {
		t.Errorf("%dx%d: title needs %d rows, has %d", w, h, len(titleLines), l.Title.InnerH())
	}
	for _, line := range titleLines {
		if runeLen(line) > l.Title.InnerW() {
			t.Errorf("%dx%d: title line %d runes > %d columns", w, h, runeLen(line), l.Title.InnerW())
		}
	}

	keyLines := []string{keyText}
	if l.KeyH == 3 {
		keyLines = []string{" ENTER=Select", " BACKSPACE=Exit"}
	}
	if len(keyLines) > l.Key.InnerH() {
		t.Errorf("%dx%d: navigation needs %d rows, has %d", w, h, len(keyLines), l.Key.InnerH())
	}
	for _, line := range keyLines {
		if runeLen(line) > l.Key.InnerW() {
			t.Errorf("%dx%d: navigation line %d runes > %d columns", w, h, runeLen(line), l.Key.InnerW())
		}
	}
}

func TestComputeLayoutUnsupported(t *testing.T) {
	cases := [][2]int{
		{0, 0},
		{1, 1},
		{10, 10},
		{40, 12},
		{80, 24},
		{99, 29},
		{100, 28},
		{107, 28},
	}
	for _, c := range cases {
		if _, ok := ComputeLayout(c[0], c[1]); ok {
			t.Errorf("ComputeLayout(%d,%d) = supported, want unsupported", c[0], c[1])
		}
	}
}

func TestComputeLayoutMinimum(t *testing.T) {
	l, ok := ComputeLayout(MinTerminalWidth, MinTerminalHeight)
	if !ok {
		t.Fatalf("ComputeLayout(%d,%d) = unsupported", MinTerminalWidth, MinTerminalHeight)
	}

	want := []struct {
		name string
		got  Rect
		want Rect
	}{
		{"title", l.Title, Rect{2, 1, 93, 4}},
		{"blink", l.Blink, Rect{94, 1, 98, 3}},
		{"menu", l.Menu, Rect{2, 5, 33, 20}},
		{"key", l.Key, Rect{2, 21, 33, 24}},
		{"log", l.Log, Rect{2, 25, 98, 28}},
		{"main", l.Main, Rect{35, 5, 98, 24}},
		{"legend", l.Legend, Rect{82, 5, 98, 24}},
		{"operation", l.Operation, Rect{69, 5, 98, 24}},
	}
	for _, entry := range want {
		if entry.got != entry.want {
			t.Errorf("%s = %+v, want %+v", entry.name, entry.got, entry.want)
		}
	}

	if l.ContentTop != 5 || l.TitleH != 3 || l.KeyH != 3 {
		t.Errorf("ContentTop/TitleH/KeyH = %d/%d/%d, want 5/3/3", l.ContentTop, l.TitleH, l.KeyH)
	}
	if got, want := l.Main.InnerW(), 62; got != want {
		t.Errorf("main inner width = %d, want %d", got, want)
	}
	if got, want := l.Main.InnerH(), 18; got != want {
		t.Errorf("main inner height = %d, want %d", got, want)
	}
	if got, want := l.Legend.InnerW(), 15; got != want {
		t.Errorf("legend inner width = %d, want %d", got, want)
	}
	if got, want := l.Operation.InnerW(), 28; got != want {
		t.Errorf("operation inner width = %d, want %d", got, want)
	}

	checkLayoutInvariants(t, MinTerminalWidth, MinTerminalHeight, l)
}

func TestComputeLayoutInvariants(t *testing.T) {
	sizes := [][2]int{
		{100, 29}, {101, 29}, {102, 29}, {106, 29}, {107, 29},
		{120, 30}, {160, 45}, {240, 60}, {300, 100},
	}
	for _, size := range sizes {
		l, ok := ComputeLayout(size[0], size[1])
		if !ok {
			t.Errorf("ComputeLayout(%d,%d) = unsupported", size[0], size[1])
			continue
		}
		checkLayoutInvariants(t, size[0], size[1], l)
	}
}

func TestComputeLayoutMonotonicWidth(t *testing.T) {
	for _, h := range []int{29, 30, 45, 60, 100} {
		prev, ok := ComputeLayout(MinTerminalWidth, h)
		if !ok {
			t.Fatalf("ComputeLayout(%d,%d) = unsupported", MinTerminalWidth, h)
		}
		for w := MinTerminalWidth + 1; w <= 400; w++ {
			cur, ok := ComputeLayout(w, h)
			if !ok {
				t.Fatalf("ComputeLayout(%d,%d) = unsupported", w, h)
			}

			checks := []struct {
				name      string
				prev, cur int
			}{
				{"blink width", prev.Blink.InnerW(), cur.Blink.InnerW()},
				{"blink height", prev.Blink.InnerH(), cur.Blink.InnerH()},
				{"menu width", prev.Menu.InnerW(), cur.Menu.InnerW()},
				{"menu height", prev.Menu.InnerH(), cur.Menu.InnerH()},
				{"title width", prev.Title.InnerW(), cur.Title.InnerW()},
				{"key width", prev.Key.InnerW(), cur.Key.InnerW()},
				{"log width", prev.Log.InnerW(), cur.Log.InnerW()},
				{"log height", prev.Log.InnerH(), cur.Log.InnerH()},
				{"main width", prev.Main.InnerW(), cur.Main.InnerW()},
				{"main height", prev.Main.InnerH(), cur.Main.InnerH()},
				{"legend width", prev.Legend.InnerW(), cur.Legend.InnerW()},
				{"legend height", prev.Legend.InnerH(), cur.Legend.InnerH()},
				{"operation width", prev.Operation.InnerW(), cur.Operation.InnerW()},
				{"operation height", prev.Operation.InnerH(), cur.Operation.InnerH()},
			}
			for _, check := range checks {
				if check.cur < check.prev {
					t.Errorf("%dx%d -> %dx%d: %s shrank %d -> %d", w-1, h, w, h, check.name, check.prev, check.cur)
				}
			}
			if cur.ContentTop > prev.ContentTop {
				t.Errorf("%dx%d: content top increased", w, h)
			}
			prev = cur
		}
	}
}

func TestComputeLayoutMonotonicHeight(t *testing.T) {
	for _, w := range []int{100, 101, 106, 107, 120, 240} {
		prev, ok := ComputeLayout(w, MinTerminalHeight)
		if !ok {
			t.Fatalf("ComputeLayout(%d,%d) = unsupported", w, MinTerminalHeight)
		}
		for h := MinTerminalHeight + 1; h <= 200; h++ {
			cur, ok := ComputeLayout(w, h)
			if !ok {
				t.Fatalf("ComputeLayout(%d,%d) = unsupported", w, h)
			}

			checks := []struct {
				name      string
				prev, cur int
			}{
				{"blink height", prev.Blink.InnerH(), cur.Blink.InnerH()},
				{"menu height", prev.Menu.InnerH(), cur.Menu.InnerH()},
				{"key height", prev.Key.InnerH(), cur.Key.InnerH()},
				{"log height", prev.Log.InnerH(), cur.Log.InnerH()},
				{"main height", prev.Main.InnerH(), cur.Main.InnerH()},
				{"legend height", prev.Legend.InnerH(), cur.Legend.InnerH()},
				{"operation height", prev.Operation.InnerH(), cur.Operation.InnerH()},
			}
			for _, check := range checks {
				if check.cur < check.prev {
					t.Errorf("%dx%d -> %dx%d: %s shrank %d -> %d", w, h-1, w, h, check.name, check.prev, check.cur)
				}
			}
			prev = cur
		}
	}
}

func TestComputeLayoutResizeRegressions(t *testing.T) {
	tall1, _ := ComputeLayout(120, 54)
	tall2, _ := ComputeLayout(120, 55)
	if tall2.Log.InnerH() < tall1.Log.InnerH() {
		t.Errorf("log shrank when growing 54 -> 55 rows: %d -> %d", tall1.Log.InnerH(), tall2.Log.InnerH())
	}

	wide1, _ := ComputeLayout(149, 40)
	wide2, _ := ComputeLayout(150, 40)
	if wide2.Legend.InnerW() < wide1.Legend.InnerW() {
		t.Errorf("legend shrank when growing 149 -> 150 columns: %d -> %d", wide1.Legend.InnerW(), wide2.Legend.InnerW())
	}
}

func TestComputeLayoutExhaustiveValid(t *testing.T) {
	for w := MinTerminalWidth; w <= 400; w++ {
		for h := MinTerminalHeight; h <= 120; h++ {
			l, ok := ComputeLayout(w, h)
			if !ok {
				t.Fatalf("ComputeLayout(%d,%d) = unsupported", w, h)
			}
			checkLayoutInvariants(t, w, h, l)
		}
	}
}
