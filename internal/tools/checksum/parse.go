package checksum

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	reBSD = regexp.MustCompile(`^[A-Za-z0-9_-]+ \((.+)\) = ([0-9A-Fa-f]{16,128})$`)
	reGNU = regexp.MustCompile(`^([0-9A-Fa-f]{16,128})[ \t]+[* ](.+)$`)
)

// ParseChecksumFile picks BSD vs GNU based on content.
func ParseChecksumFile(input string) map[string]string {
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(input)))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		switch {
		case reBSD.MatchString(line):
			return parseBSD(input)
		case reGNU.MatchString(line):
			return parseGNU(input)
		default:
			// keep scanning in case the first line was junk
			continue
		}
	}

	return map[string]string{}
}

func parseGNU(input string) map[string]string {
	m := make(map[string]string)

	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(input)))
	for sc.Scan() {
		line := sc.Text()
		if mm := reGNU.FindStringSubmatch(line); mm != nil {
			hash, name := mm[1], mm[2]
			// Strip optional leading '*' (binary mode indicator)
			name = strings.TrimPrefix(name, "*")

			m[name] = hash
		}
	}

	return m
}

func parseBSD(input string) map[string]string {
	m := make(map[string]string)

	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(input)))
	for sc.Scan() {
		line := sc.Text()
		if mm := reBSD.FindStringSubmatch(line); mm != nil {
			name, hash := mm[1], mm[2]

			m[name] = hash
		}
	}

	return m
}

// InferAlgoFromHex infers the checksum algorithm from the length of the hex string.
//
//nolint:mnd // Magic numbers are fine here
func InferAlgoFromHex(h string) string {
	switch len(h) {
	case 32:
		return "md5"
	case 40:
		return "sha1"
	case 64:
		return "sha256"
	case 128:
		return "sha512"
	default:
		return "sha256"
	}
}
