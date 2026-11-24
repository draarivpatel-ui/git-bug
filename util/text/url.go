package text

import (
	"strconv"
	"strings"
)

// maxURLLength is the length of an HTTPS URL at most.
const maxURLLength = 2048

// ValidHTTPSURL checks the spelling of a stored HTTPS URL (values.md §9).
// The port may be omitted or written in decimal from 1 through 65535.
func ValidHTTPSURL(s string) bool {
	_, ok := parseHTTPSURL(s)
	return ok
}

// FetchableHTTPSURL tells if a client may request the URL (values.md §9.1): an
// HTTPS URL on port 443. Redirects are the fetcher's to check.
func FetchableHTTPSURL(s string) bool {
	port, ok := parseHTTPSURL(s)
	return ok && port == 443
}

// parseHTTPSURL checks the spelling of an HTTPS URL and returns its port, 443
// when omitted.
func parseHTTPSURL(s string) (uint64, bool) {
	// the rule allows ASCII only, so the length in bytes is the one in code points
	if len(s) > maxURLLength {
		return 0, false
	}
	rest, ok := strings.CutPrefix(s, "https://")
	if !ok {
		return 0, false
	}
	host, rest := rest, ""
	if i := strings.IndexAny(host, "/?#"); i >= 0 {
		host, rest = host[:i], host[i:]
	}
	port := uint64(443)
	host, portStr, hasPort := strings.Cut(host, ":")
	if hasPort {
		n, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil || n == 0 {
			return 0, false
		}
		port = n
	}
	return port, validHost(host) && validURLRest(rest)
}

// CleanupHTTPSURL returns the string if it is an HTTPS URL (see ValidHTTPSURL),
// or else the empty string: for a URL from elsewhere, such as a bridge, which
// can be left out but not kept as it is.
func CleanupHTTPSURL(s string) string {
	if ValidHTTPSURL(s) {
		return s
	}
	return ""
}

// validHost tells if the string is two or more labels separated by ".", each of
// 1 to 63 ASCII letters, digits and "-", the last one starting with a letter.
func validHost(host string) bool {
	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if len(label) < 1 || len(label) > 63 {
			return false
		}
		for i := 0; i < len(label); i++ {
			if !isASCIILetter(label[i]) && !isASCIIDigit(label[i]) && label[i] != '-' {
				return false
			}
		}
	}
	return isASCIILetter(labels[len(labels)-1][0])
}

// validURLRest tells if the string, empty or starting with "/", "?" or "#", holds
// only the characters allowed after the host, each "%" followed by two
// hexadecimal digits.
func validURLRest(rest string) bool {
	for i := 0; i < len(rest); i++ {
		c := rest[i]
		switch {
		case c == '%':
			if i+2 >= len(rest) || !isHexDigit(rest[i+1]) || !isHexDigit(rest[i+2]) {
				return false
			}
			i += 2
		case isASCIILetter(c), isASCIIDigit(c), strings.IndexByte("-._~:/?#[]@!$&'()*+,;=", c) >= 0:
		default:
			return false
		}
	}
	return true
}

func isASCIILetter(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }

func isASCIIDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHexDigit(c byte) bool {
	return isASCIIDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
