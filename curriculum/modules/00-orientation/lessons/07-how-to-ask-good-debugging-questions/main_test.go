package main

import "testing"

func TestScoreQuestion_AllFieldsFilled(t *testing.T) {
	q := DebugQuestion{
		Expected:    "server returns 200",
		Actual:      "server returns 503",
		Command:     "curl localhost:8080",
		Error:       "503 error",
		Environment: "Windows, Go 1.25",
	}
	score := scoreQuestion(q)
	if score != 5 {
		t.Errorf("got score %d, want 5", score)
	}
}

func TestScoreQuestion_NoFieldsFilled(t *testing.T) {
	q := DebugQuestion{}
	score := scoreQuestion(q)
	if score != 0 {
		t.Errorf("got score %d, want 0", score)
	}
}

func TestScoreQuestion_PartialFields(t *testing.T) {
	tests := []struct {
		name  string
		q     DebugQuestion
		want  int
	}{
		{
			name: "only expected and actual",
			q: DebugQuestion{
				Expected: "works",
				Actual:   "broken",
			},
			want: 2,
		},
		{
			name: "only error",
			q: DebugQuestion{
				Error: "panic",
			},
			want: 1,
		},
		{
			name: "all but environment",
			q: DebugQuestion{
				Expected: "a",
				Actual:   "b",
				Command:  "c",
				Error:    "d",
			},
			want: 4,
		},
		{
			name: "whitespace-only fields treated as empty",
			q: DebugQuestion{
				Expected: "   ",
				Actual:   "real output",
				Command:  "\t",
			},
			want: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := scoreQuestion(tc.q)
			if got != tc.want {
				t.Errorf("scoreQuestion() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestFormatField_Empty(t *testing.T) {
	if formatField("") != "(empty)" {
		t.Errorf("formatField('') should return '(empty)'")
	}
}

func TestFormatField_NonEmpty(t *testing.T) {
	if formatField("hello") != "hello" {
		t.Errorf("formatField('hello') should return 'hello'")
	}
}
