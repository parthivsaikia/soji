package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidReference = fmt.Errorf("invalid reference")
	ErrParsingFailed    = fmt.Errorf("parsing failed")
)

type Target struct {
	Provider string
	Host     string
	Path     string // Full path after host: "namespace/repo", "account/repo", or "repo"
	Extras   map[string]string
}

const (
	dockerLibrary        = "library"
	dockerRegistry       = "docker.io"
	dockerRegistryLegacy = "index.docker.io"
	dockerRegistryDNS    = "registry-1.docker.io"
)

// Provider constants for known container registry vendors.
const (
	ProviderDockerHub = "dockerhub"
	ProviderGHCR      = "ghcr"
	ProviderQuay      = "quay"
	ProviderECR       = "ecr"
	ProviderGCR       = "gcr"
	ProviderACR       = "acr"
	ProviderLocal     = "local"
	ProviderUnknown   = "unknown"
)

// hasNamespaceInPath returns true if the provider uses namespace/user/project as first path component.
// For these providers: host/namespace/repo. Otherwise: host/repo.
func hasNamespaceInPath(provider string) bool {
	switch provider {
	case ProviderDockerHub, ProviderGHCR, ProviderQuay, ProviderGCR:
		return true
	default:
		return false
	}
}

// detectProvider returns the provider name based on the registry hostname.
func detectProvider(host string) string {
	host = strings.ToLower(host)
	// Remove port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 && !strings.Contains(host[idx:], "]") {
		host = host[:idx]
	}

	switch {
	case host == dockerRegistry || host == dockerRegistryLegacy || host == dockerRegistryDNS:
		return ProviderDockerHub
	case host == "ghcr.io":
		return ProviderGHCR
	case strings.HasSuffix(host, ".azurecr.io"):
		return ProviderACR
	case strings.Contains(host, ".dkr.ecr.") && strings.HasSuffix(host, ".amazonaws.com"):
		return ProviderECR
	case host == "gcr.io" || strings.HasSuffix(host, ".gcr.io") || strings.HasSuffix(host, ".pkg.dev"):
		return ProviderGCR
	case strings.HasSuffix(host, "quay.io") || strings.HasSuffix(host, ".quay.io"):
		return ProviderQuay
	case host == "localhost" || strings.HasPrefix(host, "localhost:"):
		return ProviderLocal
	default:
		return ProviderUnknown
	}
}

var (
	hostPartS = `(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?)`
	portS     = `(?:` + regexp.QuoteMeta(`:`) + `[0-9]+)`
	ipv6PartS = `(?:[0-9a-fA-F]{1,4}:){0,7}[0-9a-fA-F]{1,4}`
	ipv6S     = `(?:` + regexp.QuoteMeta(`[`) + `(?:` +
		ipv6PartS + `|` +
		regexp.QuoteMeta(`::`) + ipv6PartS + `|` +
		ipv6PartS + regexp.QuoteMeta(`::`) + ipv6PartS + `|` +
		ipv6PartS + regexp.QuoteMeta(`::`) +
		`)` + regexp.QuoteMeta(`]`) + `)`
	localhostS  = `localhost`
	hostDomainS = `(?:` + hostPartS + `(?:(?:` + regexp.QuoteMeta(`.`) + hostPartS + `)+` + regexp.QuoteMeta(`.`) + `?|` + regexp.QuoteMeta(`.`) + `))`
	hostUpperS  = `(?:[a-zA-Z0-9]*[A-Z][a-zA-Z0-9-]*[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[A-Z][a-zA-Z0-9]*)`
	registryS   = `(?:` +
		`(?:` + hostDomainS + `|` + hostUpperS + `|` + ipv6S + `|` + localhostS + `)` + portS + `?|` +
		hostPartS + portS +
		`)`
	repoPartS = `[a-z0-9]+(?:(?:\.|_|__|-+)[a-z0-9]+)*`
	tagS      = `[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}`
	digestS   = `[A-Za-z][A-Za-z0-9]*(?:[-_+.][A-Za-z][A-Za-z0-9]*)*[:][[:xdigit:]]{32,}`
	schemeRE  = regexp.MustCompile(`^([a-z]+)://(.+)$`)
	registryRE = regexp.MustCompile(`^(` + registryS + `)$`)
	refRE     = regexp.MustCompile(`^(?:(` + registryS + `)` + regexp.QuoteMeta(`/`) + `)?` +
		`(` + repoPartS + `(?:` + regexp.QuoteMeta(`/`) + repoPartS + `)*)` +
		`(?:` + regexp.QuoteMeta(`:`) + `(` + tagS + `))?` +
		`(?:` + regexp.QuoteMeta(`@`) + `(` + digestS + `))?$`)
)

// Parse parses a container image reference URL into a Target struct.
// Supports formats like:
//   - docker.io/library/ubuntu:latest
//   - registry.example.com/repo:tag@sha256:digest
//   - ghcr.io/user/repo:tag
func Parse(url string) (Target, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return Target{}, fmt.Errorf("%w: empty reference", ErrInvalidReference)
	}

	scheme := ""
	tail := url
	matchScheme := schemeRE.FindStringSubmatch(url)
	if len(matchScheme) == 3 {
		scheme = matchScheme[1]
		tail = matchScheme[2]
	}

	ret := Target{
		Extras: make(map[string]string),
	}

	switch scheme {
	case "":
		matchRef := refRE.FindStringSubmatch(tail)
		if len(matchRef) < 5 {
			if refRE.FindStringSubmatch(strings.ToLower(tail)) != nil {
				return Target{}, fmt.Errorf("%w: \"%s\", repo must be lowercase", ErrInvalidReference, tail)
			}
			return Target{}, fmt.Errorf("%w: \"%s\"", ErrInvalidReference, tail)
		}
		ret.Host = matchRef[1]
		ret.Path = matchRef[2]
		if matchRef[3] != "" {
			ret.Extras["tag"] = matchRef[3]
		}
		if matchRef[4] != "" {
			ret.Extras["digest"] = matchRef[4]
		}

		// A valid reference must have a registry prefix (host:port/...)
		// Exception: "localhost" as first path component (treated as registry)
		if ret.Host == "" {
			repoPath := strings.Split(ret.Path, "/")
			if len(repoPath) > 0 && repoPath[0] == "localhost" {
				ret.Host = repoPath[0]
				ret.Path = strings.Join(repoPath[1:], "/")
			} else {
				return Target{}, fmt.Errorf("%w: \"%s\", missing registry host", ErrInvalidReference, tail)
			}
		}

		// Normalize docker registry hostnames
		switch ret.Host {
		case "", dockerRegistryDNS, dockerRegistryLegacy:
			ret.Host = dockerRegistry
		}

		// Detect provider from hostname
		ret.Provider = detectProvider(ret.Host)

		// Default tag to latest if neither tag nor digest specified
		if _, hasTag := ret.Extras["tag"]; !hasTag {
			if _, hasDigest := ret.Extras["digest"]; !hasDigest {
				ret.Extras["tag"] = "latest"
			}
		}

		// For non-namespace providers, path (repo) is required
		if !hasNamespaceInPath(ret.Provider) && ret.Path == "" {
			return Target{}, fmt.Errorf("%w: \"%s\", missing repository path", ErrInvalidReference, tail)
		}

	default:
		return Target{}, fmt.Errorf("%w: unknown scheme \"%s\" in \"%s\"", ErrInvalidReference, scheme, url)
	}

	return ret, nil
}

// ParseHost parses a registry hostname into a Target.
// Returns a Target with just the host populated.
func ParseHost(url string) (Target, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return Target{}, fmt.Errorf("%w: empty reference", ErrInvalidReference)
	}

	scheme := ""
	tail := url
	matchScheme := schemeRE.FindStringSubmatch(url)
	if len(matchScheme) == 3 {
		scheme = matchScheme[1]
		tail = matchScheme[2]
	}

	ret := Target{
		Extras: make(map[string]string),
	}

	switch scheme {
	case "":
		matchReg := registryRE.FindStringSubmatch(tail)
		if len(matchReg) < 2 {
			return Target{}, fmt.Errorf("%w: \"%s\"", ErrParsingFailed, tail)
		}
		ret.Host = matchReg[1]
		if ret.Host == "" {
			return Target{}, fmt.Errorf("%w: \"%s\"", ErrParsingFailed, tail)
		}
		ret.Provider = detectProvider(ret.Host)

	default:
		return Target{}, fmt.Errorf("%w: unknown scheme \"%s\" in \"%s\"", ErrParsingFailed, scheme, url)
	}

	return ret, nil
}

// String returns the canonical string representation of the Target.
func (t Target) String() string {
	var sb strings.Builder

	if t.Host != "" {
		sb.WriteString(t.Host)
		sb.WriteString("/")
	}
	if t.Path == "" {
		return ""
	}
	sb.WriteString(t.Path)
	if tag, ok := t.Extras["tag"]; ok {
		sb.WriteString(":")
		sb.WriteString(tag)
	}
	if digest, ok := t.Extras["digest"]; ok {
		sb.WriteString("@")
		sb.WriteString(digest)
	}

	return sb.String()
}

// WithDigest returns a new Target with the digest set.
// The tag is NOT cleared.
func (t Target) WithDigest(digest string) Target {
	t.Extras["digest"] = digest
	return t
}

// WithTag returns a new Target with the tag set.
// The digest is cleared.
func (t Target) WithTag(tag string) Target {
	t.Extras["tag"] = tag
	delete(t.Extras, "digest")
	return t
}

// GetNamespace returns the namespace/account/user part of the path for providers that use it.
// Returns empty string for providers without namespace structure.
func (t Target) GetNamespace() string {
	if !hasNamespaceInPath(t.Provider) {
		return ""
	}
	parts := strings.Split(t.Path, "/")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// GetRepo returns the repository name (without namespace) for providers that use namespace.
// For non-namespace providers, returns the full Path.
func (t Target) GetRepo() string {
	if !hasNamespaceInPath(t.Provider) {
		return t.Path
	}
	parts := strings.Split(t.Path, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[1:], "/")
	}
	return ""
}

// IsSet returns true if the Target has all required fields for a complete reference.
// For namespace providers, namespace-only (without repo) is valid for listing repos.
// For specific repos, tag or digest is required.
func (t Target) IsSet() bool {
	if !t.IsSetRepo() {
		return false
	}
	// If we have a specific repo (GetRepo() != ""), require tag or digest
	if t.GetRepo() != "" {
		if _, hasTag := t.Extras["tag"]; !hasTag {
			if _, hasDigest := t.Extras["digest"]; !hasDigest {
				return false
			}
		}
	}
	// Namespace-only (no repo) is valid for listing
	return true
}

// IsSetRepo returns true when the Target includes values for a specific repository
// or a namespace (for providers that use namespace/repo structure).
func (t Target) IsSetRepo() bool {
	if t.Host == "" {
		return false
	}
	if hasNamespaceInPath(t.Provider) {
		// For namespace providers, path must have at least namespace
		return t.Path != ""
	}
	// For non-namespace providers, path (repo) is required
	return t.Path != ""
}

// IsZero returns true if Target is unset.
func (t Target) IsZero() bool {
	return t.Provider == "" && t.Host == "" && t.Path == "" && len(t.Extras) == 0
}