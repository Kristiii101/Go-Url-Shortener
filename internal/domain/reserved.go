package domain

import "strings"

var reserved = map[string]struct{}{
	"api":         {},
	"v1":          {},
	"healthz":     {},
	"readyz":      {},
	"metrics":     {},
	"admin":       {},
	"docs":        {},
	"robots.txt":  {},
	"favicon.ico": {},
	"sitemap.xml": {},
	"static":      {},
	"assets":      {},
	"app":         {},
}

func IsReserved(key string) bool {
	_, ok := reserved[strings.ToLower(strings.TrimSpace(key))]
	return ok
}
