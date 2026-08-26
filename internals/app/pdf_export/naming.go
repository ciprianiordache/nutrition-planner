package pdf_export

import (
	"regexp"
	"strings"
)

// invalidPathChars matches characters that are illegal (or awkward) in a
// file/directory name on at least one of Windows/macOS/Linux - since
// exports get written directly to the user's filesystem (e.g. Desktop),
// we sanitize instead of trusting client/plan names as-is.
var invalidPathChars = regexp.MustCompile(`[\\/:*?"<>|]+`)

// sanitizePathSegment turns an arbitrary string (a client name, a plan
// title) into something safe to use as a single path segment: illegal
// characters become "-", and leading/trailing dots or spaces (which
// Windows rejects) are trimmed. Falls back to fallback if the result would
// otherwise be empty.
func sanitizePathSegment(s, fallback string) string {
	s = strings.TrimSpace(s)
	s = invalidPathChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, ". ")
	if s == "" {
		return fallback
	}
	return s
}
