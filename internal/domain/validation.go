package domain

import (
	"net/url"
	"regexp"
	"strings"
)

var aliasRe = regexp.MustCompile(`^[A-Za-z0-9_-]{3,32}$`)

func ValidateAlias(alias string) bool {
	return aliasRe.MatchString(alias)
}

func CanonicalizeURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", ErrInvalidURL
	}
	if u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", ErrInvalidURL
	}
	if u.Host == "" {
		return "", ErrInvalidURL
	}

	// lower-case host
	u.Host = strings.ToLower(u.Host)
	// strip fragment
	u.Fragment = ""

	// trim default ports if present
	port := u.Port()
	if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
		u.Host = u.Hostname() // removes the :port
	}

	return u.String(), nil
}
