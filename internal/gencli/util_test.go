package gencli

import "testing"

func TestTitle(t *testing.T) {
	testCases := []struct {
		name, input, want string
	}{
		{
			name:  "simple",
			input: "title",
			want:  "Title",
		},
		{
			name:  "underscores",
			input: "title_with_underscores",
			want:  "TitleWithUnderscores",
		},
		{
			name:  "dots",
			input: "title.with.dots",
			want:  "Title.With.Dots",
		},
		{
			name:  "dots and underscores",
			input: "title.with_dots.and_underscores",
			want:  "Title.WithDots.AndUnderscores",
		},
	}
	for _, tst := range testCases {
		t.Run(tst.name, func(t *testing.T) {
			if got := title(tst.input); got != tst.want {
				t.Errorf("title(%s): got %s, want %s", tst.input, got, tst.want)
			}
		})
	}
}
