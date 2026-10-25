package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		input     string
		provider  string
		host      string
		path      string // full path after host
		tag       string
		digest    string
		wantErr   bool
	}{
		// Docker Hub - namespace providers (path = namespace/repo or namespace)
		{
			input:    "docker.io/library/ubuntu:latest",
			provider: ProviderDockerHub,
			host:     "docker.io",
			path:     "library/ubuntu",
			tag:      "latest",
		},
		{
			input:    "docker.io/ubuntu:20.04",
			provider: ProviderDockerHub,
			host:     "docker.io",
			path:     "ubuntu", // namespace only
			tag:      "20.04",
		},
		{
			input:    "docker.io/parthiv171/myapp:v1",
			provider: ProviderDockerHub,
			host:     "docker.io",
			path:     "parthiv171/myapp",
			tag:      "v1",
		},
		{
			input:    "index.docker.io/library/alpine:latest",
			provider: ProviderDockerHub,
			host:     "docker.io",
			path:     "library/alpine",
			tag:      "latest",
		},
		{
			input:    "registry-1.docker.io/library/nginx:latest",
			provider: ProviderDockerHub,
			host:     "docker.io",
			path:     "library/nginx",
			tag:      "latest",
		},
		{
			input:    "docker.io/alpine",
			provider: ProviderDockerHub,
			host:     "docker.io",
			path:     "alpine", // namespace only
			tag:      "latest",
		},

		// GHCR - namespace providers
		{
			input:    "ghcr.io/user/repo:tag",
			provider: ProviderGHCR,
			host:     "ghcr.io",
			path:     "user/repo",
			tag:      "tag",
		},
		{
			input:    "ghcr.io/myorg",
			provider: ProviderGHCR,
			host:     "ghcr.io",
			path:     "myorg", // namespace only
			tag:      "latest",
		},

		// Quay - namespace providers
		{
			input:    "quay.io/user/repo:tag",
			provider: ProviderQuay,
			host:     "quay.io",
			path:     "user/repo",
			tag:      "tag",
		},
		{
			input:    "quay.io/myorg",
			provider: ProviderQuay,
			host:     "quay.io",
			path:     "myorg", // namespace only
			tag:      "latest",
		},

		// GCR - namespace providers
		{
			input:    "gcr.io/project/repo:tag",
			provider: ProviderGCR,
			host:     "gcr.io",
			path:     "project/repo",
			tag:      "tag",
		},
		{
			input:    "us-docker.pkg.dev/project/repo:tag",
			provider: ProviderGCR,
			host:     "us-docker.pkg.dev",
			path:     "project/repo",
			tag:      "tag",
		},
		{
			input:    "gcr.io/myproject",
			provider: ProviderGCR,
			host:     "gcr.io",
			path:     "myproject", // namespace only
			tag:      "latest",
		},

		// Non-namespace providers (path = repo directly)
		{
			input:    "registry.example.com:5000/repo:tag@sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			provider: ProviderUnknown,
			host:     "registry.example.com:5000",
			path:     "repo",
			tag:      "tag",
			digest:   "sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		},
		{
			input:    "localhost:5000/myapp:v1",
			provider: ProviderLocal,
			host:     "localhost:5000",
			path:     "myapp",
			tag:      "v1",
		},
		{
			input:    "myregistry.com:5000/path/to/repo",
			provider: ProviderUnknown,
			host:     "myregistry.com:5000",
			path:     "path/to/repo",
			tag:      "latest",
		},
		{
			input:    "localhost/myrepo:tag",
			provider: ProviderLocal,
			host:     "localhost",
			path:     "myrepo",
			tag:      "tag",
		},
		{
			input:    "myacr.azurecr.io/repo:tag",
			provider: ProviderACR,
			host:     "myacr.azurecr.io",
			path:     "repo",
			tag:      "tag",
		},
		{
			input:    "123456789012.dkr.ecr.us-east-1.amazonaws.com/repo:tag",
			provider: ProviderECR,
			host:     "123456789012.dkr.ecr.us-east-1.amazonaws.com",
			path:     "repo",
			tag:      "tag",
		},

		// Error cases
		{
			input:   "",
			wantErr: true,
		},
		{
			input:   "invalid@reference",
			wantErr: true,
		},
		{
			input:   "unknown://scheme/path",
			wantErr: true,
		},
		// Non-namespace providers require a path
		{
			input:   "myacr.azurecr.io",
			wantErr: true,
		},
		{
			input:   "123456789012.dkr.ecr.us-east-1.amazonaws.com",
			wantErr: true,
		},
		{
			input:   "localhost:5000",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tc.input, err)
				return
			}
			if result.Provider != tc.provider {
				t.Errorf("provider: got %q, want %q", result.Provider, tc.provider)
			}
			if result.Host != tc.host {
				t.Errorf("host: got %q, want %q", result.Host, tc.host)
			}
			if result.Path != tc.path {
				t.Errorf("path: got %q, want %q", result.Path, tc.path)
			}
			if tag, ok := result.Extras["tag"]; ok {
				if tag != tc.tag {
					t.Errorf("tag: got %q, want %q", tag, tc.tag)
				}
			} else if tc.tag != "" {
				t.Errorf("tag: missing, want %q", tc.tag)
			}
			if digest, ok := result.Extras["digest"]; ok {
				if digest != tc.digest {
					t.Errorf("digest: got %q, want %q", digest, tc.digest)
				}
			} else if tc.digest != "" {
				t.Errorf("digest: missing, want %q", tc.digest)
			}
		})
	}
}

func TestParseHost(t *testing.T) {
	testCases := []struct {
		input    string
		provider string
		host     string
		path     string
		wantErr  bool
	}{
		{
			input:    "docker.io",
			provider: ProviderDockerHub,
			host:     "docker.io",
		},
		{
			input:    "registry.example.com:5000",
			provider: ProviderUnknown,
			host:     "registry.example.com:5000",
		},
		{
			input:    "localhost:5000",
			provider: ProviderLocal,
			host:     "localhost:5000",
		},
		{
			input:    "ghcr.io",
			provider: ProviderGHCR,
			host:     "ghcr.io",
		},
		{
			input:    "quay.io",
			provider: ProviderQuay,
			host:     "quay.io",
		},
		{
			input:    "myacr.azurecr.io",
			provider: ProviderACR,
			host:     "myacr.azurecr.io",
		},
		{
			input:   "",
			wantErr: true,
		},
		{
			input:   "invalid@host",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := ParseHost(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tc.input, err)
				return
			}
			if result.Provider != tc.provider {
				t.Errorf("provider: got %q, want %q", result.Provider, tc.provider)
			}
			if result.Host != tc.host {
				t.Errorf("host: got %q, want %q", result.Host, tc.host)
			}
			if result.Path != tc.path {
				t.Errorf("path: got %q, want %q", result.Path, tc.path)
			}
		})
	}
}

func TestTargetString(t *testing.T) {
	testCases := []struct {
		target Target
		want   string
	}{
		{
			target: Target{
				Provider: ProviderDockerHub,
				Host:     "docker.io",
				Path:     "library/ubuntu",
				Extras:   map[string]string{"tag": "latest"},
			},
			want: "docker.io/library/ubuntu:latest",
		},
		{
			target: Target{
				Provider: ProviderGHCR,
				Host:     "ghcr.io",
				Path:     "user/repo",
				Extras:   map[string]string{"tag": "v1", "digest": "sha256:abc123"},
			},
			want: "ghcr.io/user/repo:v1@sha256:abc123",
		},
		{
			target: Target{
				Provider: ProviderDockerHub,
				Host:     "docker.io",
				Path:     "myorg", // namespace only
			},
			want: "docker.io/myorg",
		},
		{
			target: Target{
				Provider: ProviderUnknown,
				Host:     "docker.io",
				Path:     "library/alpine",
			},
			want: "docker.io/library/alpine",
		},
		{
			target: Target{
				Provider: ProviderECR,
				Host:     "123.dkr.ecr.us-east-1.amazonaws.com",
				Path:     "repo",
				Extras:   map[string]string{"tag": "latest"},
			},
			want: "123.dkr.ecr.us-east-1.amazonaws.com/repo:latest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			got := tc.target.String()
			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTargetMethods(t *testing.T) {
	t.Run("WithDigest", func(t *testing.T) {
		target := Target{
			Provider: ProviderDockerHub,
			Host:     "docker.io",
			Path:     "library/ubuntu",
			Extras:   map[string]string{"tag": "latest"},
		}
		result := target.WithDigest("sha256:abc123")
		if result.Extras["digest"] != "sha256:abc123" {
			t.Errorf("WithDigest: digest not set")
		}
		if result.Extras["tag"] != "latest" {
			t.Errorf("WithDigest: tag should not be cleared")
		}
	})

	t.Run("WithTag", func(t *testing.T) {
		target := Target{
			Provider: ProviderDockerHub,
			Host:     "docker.io",
			Path:     "library/ubuntu",
			Extras:   map[string]string{"digest": "sha256:abc123"},
		}
		result := target.WithTag("v1.0")
		if result.Extras["tag"] != "v1.0" {
			t.Errorf("WithTag: tag not set")
		}
		if _, ok := result.Extras["digest"]; ok {
			t.Errorf("WithTag: digest should be cleared")
		}
	})

	t.Run("GetNamespace", func(t *testing.T) {
		tests := []struct {
			target    Target
			want      string
			desc      string
		}{
			{Target{Provider: ProviderDockerHub, Path: "library/ubuntu"}, "library", "dockerhub with namespace/repo"},
			{Target{Provider: ProviderDockerHub, Path: "myorg"}, "myorg", "dockerhub namespace only"},
			{Target{Provider: ProviderGHCR, Path: "user/repo"}, "user", "ghcr with namespace/repo"},
			{Target{Provider: ProviderGHCR, Path: "myorg"}, "myorg", "ghcr namespace only"},
			{Target{Provider: ProviderECR, Path: "repo"}, "", "ecr no namespace"},
			{Target{Provider: ProviderACR, Path: "repo"}, "", "acr no namespace"},
			{Target{Provider: ProviderLocal, Path: "repo"}, "", "local no namespace"},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got := tt.target.GetNamespace()
				if got != tt.want {
					t.Errorf("GetNamespace() = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("GetRepo", func(t *testing.T) {
		tests := []struct {
			target    Target
			want      string
			desc      string
		}{
			{Target{Provider: ProviderDockerHub, Path: "library/ubuntu"}, "ubuntu", "dockerhub with namespace/repo"},
			{Target{Provider: ProviderDockerHub, Path: "myorg"}, "", "dockerhub namespace only (no repo)"},
			{Target{Provider: ProviderGHCR, Path: "user/repo"}, "repo", "ghcr with namespace/repo"},
			{Target{Provider: ProviderGHCR, Path: "myorg"}, "", "ghcr namespace only (no repo)"},
			{Target{Provider: ProviderECR, Path: "repo"}, "repo", "ecr no namespace"},
			{Target{Provider: ProviderACR, Path: "repo"}, "repo", "acr no namespace"},
			{Target{Provider: ProviderLocal, Path: "repo"}, "repo", "local no namespace"},
			{Target{Provider: ProviderGCR, Path: "project/repo/path"}, "repo/path", "gcr with nested path"},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got := tt.target.GetRepo()
				if got != tt.want {
					t.Errorf("GetRepo() = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("IsSet", func(t *testing.T) {
		tests := []struct {
			target Target
			want   bool
			desc   string
		}{
			// Namespace providers: namespace-only is valid (for listing)
			{Target{Provider: ProviderDockerHub, Host: "docker.io", Path: "myorg", Extras: map[string]string{}}, true, "dockerhub namespace only"},
			{Target{Provider: ProviderGHCR, Host: "ghcr.io", Path: "myorg", Extras: map[string]string{}}, true, "ghcr namespace only"},
			// Namespace providers: namespace/repo requires tag or digest
			{Target{Provider: ProviderDockerHub, Host: "docker.io", Path: "library/ubuntu", Extras: map[string]string{"tag": "latest"}}, true, "dockerhub namespace/repo with tag"},
			{Target{Provider: ProviderDockerHub, Host: "docker.io", Path: "library/ubuntu", Extras: map[string]string{"digest": "sha256:abc"}}, true, "dockerhub namespace/repo with digest"},
			{Target{Provider: ProviderDockerHub, Host: "docker.io", Path: "library/ubuntu", Extras: map[string]string{}}, false, "dockerhub namespace/repo no tag/digest"},
			// Non-namespace providers: path required + tag/digest
			{Target{Provider: ProviderECR, Host: "123.dkr.ecr.us-east-1.amazonaws.com", Path: "repo", Extras: map[string]string{"tag": "latest"}}, true, "ecr with tag"},
			{Target{Provider: ProviderECR, Host: "123.dkr.ecr.us-east-1.amazonaws.com", Path: "repo", Extras: map[string]string{"digest": "sha256:abc"}}, true, "ecr with digest"},
			{Target{Provider: ProviderECR, Host: "123.dkr.ecr.us-east-1.amazonaws.com", Path: "repo", Extras: map[string]string{}}, false, "ecr no tag/digest"},
			{Target{Provider: ProviderECR, Host: "123.dkr.ecr.us-east-1.amazonaws.com", Path: "", Extras: map[string]string{}}, false, "ecr no path"},
			{Target{}, false, "empty"},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got := tt.target.IsSet()
				if got != tt.want {
					t.Errorf("IsSet() = %v, want %v for %+v", got, tt.want, tt.target)
				}
			})
		}
	})

	t.Run("IsSetRepo", func(t *testing.T) {
		tests := []struct {
			target Target
			want   bool
			desc   string
		}{
			// Namespace providers: namespace counts as repo set
			{Target{Provider: ProviderDockerHub, Host: "docker.io", Path: "myorg"}, true, "dockerhub namespace"},
			{Target{Provider: ProviderDockerHub, Host: "docker.io", Path: "library/ubuntu"}, true, "dockerhub namespace/repo"},
			{Target{Provider: ProviderGHCR, Host: "ghcr.io", Path: "myorg"}, true, "ghcr namespace"},
			// Non-namespace providers: path required
			{Target{Provider: ProviderECR, Host: "123.dkr.ecr.us-east-1.amazonaws.com", Path: "repo"}, true, "ecr with repo"},
			{Target{Provider: ProviderECR, Host: "123.dkr.ecr.us-east-1.amazonaws.com", Path: ""}, false, "ecr no path"},
			{Target{}, false, "empty"},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got := tt.target.IsSetRepo()
				if got != tt.want {
					t.Errorf("IsSetRepo() = %v, want %v for %+v", got, tt.want, tt.target)
				}
			})
		}
	})

	t.Run("IsZero", func(t *testing.T) {
		empty := Target{}
		if !empty.IsZero() {
			t.Errorf("empty Target should be zero")
		}
		nonEmpty := Target{Provider: ProviderDockerHub}
		if nonEmpty.IsZero() {
			t.Errorf("non-empty Target should not be zero")
		}
	})
}

func TestDetectProvider(t *testing.T) {
	testCases := []struct {
		host     string
		provider string
	}{
		{"docker.io", ProviderDockerHub},
		{"index.docker.io", ProviderDockerHub},
		{"registry-1.docker.io", ProviderDockerHub},
		{"docker.io:5000", ProviderDockerHub},
		{"ghcr.io", ProviderGHCR},
		{"ghcr.io:443", ProviderGHCR},
		{"myacr.azurecr.io", ProviderACR},
		{"myacr.azurecr.io:443", ProviderACR},
		{"123456789012.dkr.ecr.us-east-1.amazonaws.com", ProviderECR},
		{"123456789012.dkr.ecr.us-east-1.amazonaws.com:443", ProviderECR},
		{"gcr.io", ProviderGCR},
		{"us.gcr.io", ProviderGCR},
		{"eu.gcr.io", ProviderGCR},
		{"us-docker.pkg.dev", ProviderGCR},
		{"europe-docker.pkg.dev", ProviderGCR},
		{"quay.io", ProviderQuay},
		{"myquay.quay.io", ProviderQuay},
		{"localhost", ProviderLocal},
		{"localhost:5000", ProviderLocal},
		{"myregistry.com", ProviderUnknown},
		{"myregistry.com:5000", ProviderUnknown},
	}

	for _, tc := range testCases {
		t.Run(tc.host, func(t *testing.T) {
			got := detectProvider(tc.host)
			if got != tc.provider {
				t.Errorf("detectProvider(%q) = %q, want %q", tc.host, got, tc.provider)
			}
		})
	}
}