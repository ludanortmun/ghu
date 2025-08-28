package internal

import (
	"testing"
)

func TestInferTargetFromUrl(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		want        GitHubTarget
		expectError bool
	}{
		{
			name: "basic repository URL",
			url:  "https://github.com/owner/repo",
			want: GitHubTarget{
				Owner:      "owner",
				Repository: "repo",
			},
			expectError: false,
		},
		{
			name: "repository URL with branch",
			url:  "https://github.com/owner/repo/tree/main",
			want: GitHubTarget{
				Owner:      "owner",
				Repository: "repo",
				Ref:        "main",
			},
			expectError: false,
		},
		{
			name: "repository URL with branch and directory",
			url:  "https://github.com/owner/repo/tree/main/path/to/dir",
			want: GitHubTarget{
				Owner:      "owner",
				Repository: "repo",
				Ref:        "main",
				Directory:  "path/to/dir",
			},
			expectError: false,
		},
		{
			name: "repository URL with commit hash",
			url:  "https://github.com/owner/repo/tree/a1b2c3d",
			want: GitHubTarget{
				Owner:      "owner",
				Repository: "repo",
				Ref:        "a1b2c3d",
			},
			expectError: false,
		},
		{
			name:        "invalid URL - not GitHub",
			url:         "https://gitlab.com/owner/repo",
			expectError: true,
		},
		{
			name:        "invalid URL - missing repository",
			url:         "https://github.com/owner",
			expectError: true,
		},
		{
			name:        "invalid URL - invalid tree structure",
			url:         "https://github.com/owner/repo/invalid/main",
			expectError: true,
		},
		{
			name: "URL with encoded characters",
			url:  "https://github.com/owner/repo/tree/main/path%20with%20spaces",
			want: GitHubTarget{
				Owner:      "owner",
				Repository: "repo",
				Ref:        "main",
				Directory:  "path with spaces",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InferTargetFromUrl(tt.url)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got.Owner != tt.want.Owner {
				t.Errorf("Owner = %v, want %v", got.Owner, tt.want.Owner)
			}
			if got.Repository != tt.want.Repository {
				t.Errorf("Repository = %v, want %v", got.Repository, tt.want.Repository)
			}
			if got.Ref != tt.want.Ref {
				t.Errorf("Ref = %v, want %v", got.Ref, tt.want.Ref)
			}
			if got.Directory != tt.want.Directory {
				t.Errorf("Directory = %v, want %v", got.Directory, tt.want.Directory)
			}
		})
	}
}
