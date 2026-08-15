package validate

import "strings"

// Slugify converts the provided value into a slug: it lowercases the
// value, replaces any sequence of characters that are not letters or
// digits with a single dash, and trims leading and trailing dashes.
func Slugify(value string) string {
	var sb strings.Builder
	sb.Grow(len(value))

	dash := false
	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			if dash && sb.Len() > 0 {
				sb.WriteByte('-')
			}
			sb.WriteRune(r)
			dash = false
			continue
		}
		dash = true
	}

	return sb.String()
}
