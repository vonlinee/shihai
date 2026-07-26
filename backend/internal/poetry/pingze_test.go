package poetry

import "testing"

func TestRecognizePingzeTextUsesPinyinTone(t *testing.T) {
	got := RecognizePingzeText("我爱中国，")
	want := "仄仄平平，"

	if got != want {
		t.Fatalf("RecognizePingzeText() = %q, want %q", got, want)
	}
}

func TestRecognizePingzeTextKeepsPunctuation(t *testing.T) {
	got := RecognizePingzeText("山，水。")
	want := "平，仄。"

	if got != want {
		t.Fatalf("RecognizePingzeText() = %q, want %q", got, want)
	}
}

func TestRecognizePingzeTextKeepsNonHanCharacters(t *testing.T) {
	got := RecognizePingzeText("山 水 123 -")
	want := "平 仄 123 -"

	if got != want {
		t.Fatalf("RecognizePingzeText() = %q, want %q", got, want)
	}
}

func TestPingzeFromPinyinMapsToneToPingze(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "first tone", in: "zhong1", want: PingzePing},
		{name: "second tone", in: "guo2", want: PingzePing},
		{name: "third tone", in: "wo3", want: PingzeZe},
		{name: "fourth tone", in: "ai4", want: PingzeZe},
		{name: "unknown tone", in: "de", want: PingzeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PingzeFromPinyin(tt.in); got != tt.want {
				t.Fatalf("PingzeFromPinyin(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPreparePingzeLinesForSaveAutoRecognizesWhenBlank(t *testing.T) {
	got, err := PreparePingzeLinesForSave([]string{"我爱你"}, []string{""})

	if err != nil {
		t.Fatalf("PreparePingzeLinesForSave error = %v, want nil", err)
	}
	assertStringSliceEqual(t, got, []string{"仄仄仄"})
}

func TestNormalizePingzeLinesAllowsUnknownAndNonHanMark(t *testing.T) {
	got, err := NormalizePingzeLines([]string{"abc"}, []string{"平， 仄? -"})

	if err != nil {
		t.Fatalf("NormalizePingzeLines error = %v, want nil", err)
	}
	assertStringSliceEqual(t, got, []string{"平， 仄? -"})
}

func TestNormalizePingzeLinesRejectsInvalidMark(t *testing.T) {
	_, err := NormalizePingzeLines([]string{"abc"}, []string{"平中仄"})

	if err == nil {
		t.Fatal("NormalizePingzeLines error = nil, want error")
	}
}

func assertStringSliceEqual(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("slice[%d] = %q, want %q; got %v", i, got[i], want[i], got)
		}
	}
}
