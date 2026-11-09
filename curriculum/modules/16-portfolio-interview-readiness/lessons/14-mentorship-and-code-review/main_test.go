package main

import "testing"

func TestAnalyzeReview_Basic(t *testing.T) {
	r := Review{
		Reviewer: "Alice",
		Comments: []ReviewComment{
			{Text: "Good code", Category: "positive"},
			{Text: "Consider refactoring", Category: "suggestion"},
			{Text: "This error is swallowed", Category: "critical"},
			{Text: "Why is this needed?", Category: "question"},
		},
	}
	a := AnalyzeReview(r)
	if a.TotalComments != 4 {
		t.Errorf("TotalComments = %d, want 4", a.TotalComments)
	}
	if a.PositiveCount != 1 {
		t.Errorf("PositiveCount = %d, want 1", a.PositiveCount)
	}
	if a.CriticalCount != 1 {
		t.Errorf("CriticalCount = %d, want 1", a.CriticalCount)
	}
}

func TestAnalyzeReview_Empty(t *testing.T) {
	r := Review{Reviewer: "Bob"}
	a := AnalyzeReview(r)
	if a.TotalComments != 0 {
		t.Errorf("TotalComments = %d, want 0", a.TotalComments)
	}
	if a.Sentiment != "neutral" {
		t.Errorf("Sentiment = %s, want neutral", a.Sentiment)
	}
}

func TestAnalyzeReview_Sentiment(t *testing.T) {
	tests := []struct {
		name     string
		comments []ReviewComment
		want     string
	}{
		{
			name: "mostly positive",
			comments: []ReviewComment{
				{Text: "Great", Category: "positive"},
				{Text: "Excellent", Category: "positive"},
				{Text: "Consider fix", Category: "critical"},
			},
			want: "positive",
		},
		{
			name: "mostly critical",
			comments: []ReviewComment{
				{Text: "Wrong", Category: "critical"},
				{Text: "Bad", Category: "critical"},
				{Text: "Nice", Category: "positive"},
			},
			want: "negative",
		},
		{
			name: "balanced",
			comments: []ReviewComment{
				{Text: "Good", Category: "positive"},
				{Text: "Fix this", Category: "critical"},
				{Text: "Consider refactoring", Category: "suggestion"},
				{Text: "Nice approach", Category: "positive"},
			},
			want: "balanced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Review{Reviewer: "Tester", Comments: tt.comments}
			a := AnalyzeReview(r)
			if a.Sentiment != tt.want {
				t.Errorf("Sentiment = %s, want %s", a.Sentiment, tt.want)
			}
		})
	}
}

func TestAnalyzeReview_Actionability(t *testing.T) {
	r := Review{
		Reviewer: "Mentor",
		Comments: []ReviewComment{
			{Text: "Consider using a buffered channel here", Category: "suggestion"},
			{Text: "This is wrong", Category: "critical"},
			{Text: "Suggest renaming to parseConfig", Category: "suggestion"},
		},
	}
	a := AnalyzeReview(r)
	if a.Actionable < 1 {
		t.Error("Actionability should be > 0 for comments with action words")
	}
}
