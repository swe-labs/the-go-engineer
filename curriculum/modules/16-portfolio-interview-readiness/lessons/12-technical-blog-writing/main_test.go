package main

import (
	"strings"
	"testing"
)

func TestAnalyzePost_Basic(t *testing.T) {
	post := BlogPost{
		Title:  "Test",
		Author: "Me",
		Sections: []BlogSection{
			{Heading: "Introduction", WordCount: 200, CodeBlocks: 0},
			{Heading: "Body", WordCount: 500, CodeBlocks: 5},
			{Heading: "Conclusion", WordCount: 100, CodeBlocks: 0},
		},
	}
	q := AnalyzePost(post)
	if q.TotalWords != 800 {
		t.Errorf("TotalWords = %d, want 800", q.TotalWords)
	}
	if !q.HasIntro {
		t.Error("HasIntro should be true")
	}
	if !q.HasConclusion {
		t.Error("HasConclusion should be true")
	}
}

func TestAnalyzePost_MissingSections(t *testing.T) {
	post := BlogPost{
		Title:    "No Sections",
		Author:   "Test",
		Sections: []BlogSection{},
	}
	q := AnalyzePost(post)
	if q.TotalWords != 0 {
		t.Errorf("TotalWords = %d, want 0", q.TotalWords)
	}
}

func TestAnalyzePost_Issues(t *testing.T) {
	tests := []struct {
		name         string
		post         BlogPost
		wantIssues   bool
		wantReadable string
	}{
		{
			name: "too short",
			post: BlogPost{
				Title:    "Short",
				Author:   "Me",
				Sections: []BlogSection{{Heading: "Intro", WordCount: 100, CodeBlocks: 0}},
			},
			wantIssues:   true,
			wantReadable: "too short",
		},
		{
			name: "good post",
			post: BlogPost{
				Title:  "Good",
				Author: "Me",
				Sections: []BlogSection{
					{Heading: "Introduction", WordCount: 200, CodeBlocks: 0},
					{Heading: "Body", WordCount: 600, CodeBlocks: 3},
					{Heading: "Conclusion", WordCount: 200, CodeBlocks: 0},
				},
			},
			wantIssues:   false,
			wantReadable: "good",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := AnalyzePost(tt.post)
			if (len(q.Issues) > 0) != tt.wantIssues {
				t.Errorf("issues = %v, want issues=%v", q.Issues, tt.wantIssues)
			}
			if q.Readability != tt.wantReadable {
				t.Errorf("Readability = %s, want %s", q.Readability, tt.wantReadable)
			}
		})
	}
}

func TestGenerateOutline(t *testing.T) {
	post := BlogPost{
		Title:  "My Post",
		Author: "Author",
		Sections: []BlogSection{
			{Heading: "Intro", WordCount: 100, CodeBlocks: 0},
			{Heading: "Body", WordCount: 300, CodeBlocks: 2},
		},
	}
	outline := GenerateOutline(post)
	if !strings.Contains(outline, "My Post") {
		t.Error("outline should contain title")
	}
	if !strings.Contains(outline, "Body") {
		t.Error("outline should contain section headings")
	}
	if !strings.Contains(outline, "~300 words") {
		t.Error("outline should contain word counts")
	}
}
