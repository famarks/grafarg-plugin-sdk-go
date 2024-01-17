package useragent

import (
	"errors"
	"regexp"
)

var (
	userAgentRegex   = regexp.MustCompile(`^Grafarg/([0-9]+\.[0-9]+\.[0-9]+(?:-[a-zA-Z0-9]+)?) \(([a-zA-Z0-9]+); ([a-zA-Z0-9]+)\)$`)
	errInvalidFormat = errors.New("invalid user agent format")
)

// UserAgent represents a Grafarg user agent.
// Its format is "Grafarg/<version> (<os>; <arch>)"
// Example: "Grafarg/7.0.0-beta1 (darwin; amd64)", "Grafarg/10.0.0 (windows; x86)"
type UserAgent struct {
	grafargVersion string
	arch           string
	os             string
}

// New creates a new UserAgent.
// The version must be a valid semver string, and the os and arch must be valid strings.
func New(grafargVersion, os, arch string) (*UserAgent, error) {
	ua := &UserAgent{
		grafargVersion: grafargVersion,
		os:             os,
		arch:           arch,
	}

	return Parse(ua.String())
}

// Parse creates a new UserAgent from a string.
func Parse(s string) (*UserAgent, error) {
	matches := userAgentRegex.FindStringSubmatch(s)
	if len(matches) != 4 {
		return nil, errInvalidFormat
	}

	return &UserAgent{
		grafargVersion: matches[1],
		os:             matches[2],
		arch:           matches[3],
	}, nil
}

func (ua *UserAgent) GrafargVersion() string {
	return ua.grafargVersion
}

func (ua *UserAgent) String() string {
	return "Grafarg/" + ua.grafargVersion + " (" + ua.os + "; " + ua.arch + ")"
}
