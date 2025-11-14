package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"net"
	"net/url"
	"strings"
)

// HashVisitor creates a hash from IP and User-Agent for visitor identification
func HashVisitor(ip, userAgent string) string {
	h := sha256.New()
	h.Write([]byte(ip))
	h.Write([]byte(userAgent))
	return hex.EncodeToString(h.Sum(nil))[:16] // Use first 16 chars
}

// HashString creates a SHA256 hash of a string
func HashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// ShortHash creates a short hash suitable for URLs
func ShortHash(s string) string {
	h := fnv.New64()
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)[:8]
}

// CanonicalizeURL normalizes a URL for consistent storage
func CanonicalizeURL(rawURL string) (string, error) {
	// Add scheme if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Validate
	if u.Host == "" {
		return "", fmt.Errorf("missing host")
	}

	// Normalize
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Fragment = "" // Remove fragment

	// Sort query params (optional)
	if u.RawQuery != "" {
		u.RawQuery = url.QueryEscape(u.Query().Encode())
	}

	return u.String(), nil
}

// ExtractIP extracts the real IP from request headers
func ExtractIP(remoteAddr string, headers map[string][]string) string {
	// Check common proxy headers
	if xff := headers["X-Forwarded-For"]; len(xff) > 0 && xff[0] != "" {
		ips := strings.Split(xff[0], ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	if xri := headers["X-Real-Ip"]; len(xri) > 0 && xri[0] != "" {
		if net.ParseIP(xri[0]) != nil {
			return xri[0]
		}
	}

	// Fall back to RemoteAddr
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}

	return remoteAddr
}

// MD5Hash creates an MD5 hash (use only for non-security purposes)
func MD5Hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
