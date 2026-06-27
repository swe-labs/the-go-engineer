package main

import (
	"testing"
)

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewRepo(t *testing.T) {
	r := NewRepo()
	if len(r.Working) != 0 {
		t.Errorf("NewRepo().Working should be empty, got %d entries", len(r.Working))
	}
	if len(r.StageNames()) != 0 {
		t.Errorf("NewRepo().StageNames() should be empty, got %d", len(r.StageNames()))
	}
	if len(r.History) != 0 {
		t.Errorf("NewRepo().History should be empty, got %d entries", len(r.History))
	}
}

func TestStage(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Repo)
		stageFile string
		wantErr   bool
		wantNames []string
	}{
		{
			name:      "stage existing file",
			setup:     func(r *Repo) { r.Working["a.txt"] = "content" },
			stageFile: "a.txt",
			wantErr:   false,
			wantNames: []string{"a.txt"},
		},
		{
			name:      "stage non-existent file returns error",
			setup:     func(r *Repo) {},
			stageFile: "b.txt",
			wantErr:   true,
			wantNames: []string{},
		},
		{
			name: "stage already staged file updates content",
			setup: func(r *Repo) {
				r.Working["c.txt"] = "v2"
				r.staging["c.txt"] = "v1"
			},
			stageFile: "c.txt",
			wantErr:   false,
			wantNames: []string{"c.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepo()
			tt.setup(r)
			err := r.Stage(tt.stageFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stage() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := r.StageNames()
			if !equalStringSlices(got, tt.wantNames) {
				t.Errorf("StageNames() = %v, want %v", got, tt.wantNames)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Repo)
		msg     string
		wantErr bool
		wantLen int
	}{
		{
			name: "commit with staged files succeeds",
			setup: func(r *Repo) {
				r.Working["f1.txt"] = "content1"
				r.staging["f1.txt"] = "content1"
			},
			msg:     "feat: add f1",
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "commit with nothing staged fails",
			setup:   func(r *Repo) {},
			msg:     "empty commit",
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "commit clears staging after success",
			setup: func(r *Repo) {
				r.Working["f2.txt"] = "content2"
				r.staging["f2.txt"] = "content2"
			},
			msg:     "feat: add f2",
			wantErr: false,
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepo()
			tt.setup(r)
			err := r.Commit(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(r.History) != tt.wantLen {
				t.Errorf("History length = %d, want %d", len(r.History), tt.wantLen)
			}
			if !tt.wantErr && len(r.History) > 0 {
				last := r.History[len(r.History)-1]
				if last.Message != tt.msg {
					t.Errorf("Last commit message = %q, want %q", last.Message, tt.msg)
				}
				// Verify staging is cleared after commit
				if len(r.StageNames()) != 0 {
					t.Errorf("staging should be empty after commit, got %d files", len(r.StageNames()))
				}
			}
		})
	}
}

func TestUnstage(t *testing.T) {
	r := NewRepo()
	r.Working["x.txt"] = "x"
	r.Stage("x.txt")
	if len(r.StageNames()) != 1 {
		t.Fatal("expected 1 staged file")
	}
	r.Unstage("x.txt")
	if len(r.StageNames()) != 0 {
		t.Errorf("expected 0 staged files after unstage, got %d", len(r.StageNames()))
	}
	// Re-stage and commit to verify state transitions
	r.Stage("x.txt")
	err := r.Commit("test: add x")
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	if len(r.History) != 1 {
		t.Errorf("expected 1 commit, got %d", len(r.History))
	}
}

func TestMultipleCommits(t *testing.T) {
	r := NewRepo()
	r.Working["a.txt"] = "a"
	r.Stage("a.txt")
	r.Commit("first: a")

	r.Working["b.txt"] = "b"
	r.Stage("b.txt")
	r.Commit("second: b")

	if len(r.History) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(r.History))
	}
	if r.History[0].Message != "first: a" {
		t.Errorf("commit 0 message = %q, want %q", r.History[0].Message, "first: a")
	}
	if r.History[1].Message != "second: b" {
		t.Errorf("commit 1 message = %q, want %q", r.History[1].Message, "second: b")
	}
}
