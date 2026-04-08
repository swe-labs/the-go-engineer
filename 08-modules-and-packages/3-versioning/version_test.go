package main

import "testing"

func TestVersionString(t *testing.T) {
	t.Parallel()

	v := Version{Major: 2, Minor: 3, Patch: 4}
	if got, want := v.String(), "v2.3.4"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestVersionIsCompatible(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  Version
		right Version
		want  bool
	}{
		{
			name:  "same major is compatible",
			left:  Version{Major: 1, Minor: 2, Patch: 0},
			right: Version{Major: 1, Minor: 9, Patch: 9},
			want:  true,
		},
		{
			name:  "different major is breaking",
			left:  Version{Major: 2, Minor: 0, Patch: 0},
			right: Version{Major: 1, Minor: 9, Patch: 9},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.left.IsCompatible(tt.right); got != tt.want {
				t.Fatalf("IsCompatible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionIsNewer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  Version
		right Version
		want  bool
	}{
		{
			name:  "major version wins first",
			left:  Version{Major: 2, Minor: 0, Patch: 0},
			right: Version{Major: 1, Minor: 99, Patch: 99},
			want:  true,
		},
		{
			name:  "minor version breaks tie",
			left:  Version{Major: 1, Minor: 2, Patch: 0},
			right: Version{Major: 1, Minor: 1, Patch: 9},
			want:  true,
		},
		{
			name:  "patch version breaks final tie",
			left:  Version{Major: 1, Minor: 1, Patch: 2},
			right: Version{Major: 1, Minor: 1, Patch: 1},
			want:  true,
		},
		{
			name:  "older version is not newer",
			left:  Version{Major: 1, Minor: 0, Patch: 0},
			right: Version{Major: 1, Minor: 0, Patch: 1},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.left.IsNewer(tt.right); got != tt.want {
				t.Fatalf("IsNewer() = %v, want %v", got, tt.want)
			}
		})
	}
}
