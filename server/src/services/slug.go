package services

import (
	"regexp"
	"strings"
)

var slugNonAlpha = regexp.MustCompile(`[^a-z0-9]+`)

func generateSlug(input string) string {
	slug := strings.ToLower(strings.TrimSpace(input))
	slug = slugNonAlpha.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "untitled"
	}
	return slug
}


