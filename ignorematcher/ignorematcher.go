package ignorematcher

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"
)

func IgnoreMatches(pattern, path string) (bool, error) {
	return regexpMatch(pattern, path)
}

func Chop(str string) string {
	if str == "" {
		return ""
	}
	return str[0 : len(str)-1]
}

// regexpMatch tries to match the logic of filepath.Match but
// does so using regexp logic. We do this so that we can expand the
// wildcard set to include other things, like "**" to mean any number
// of directories.  This means that we should be backwards compatible
// with filepath.Match(). We'll end up supporting more stuff, due to
// the fact that we're using regexp, but that's ok - it does no harm.
func regexpMatch(pattern, path string) (bool, error) {
	if strings.HasSuffix(pattern, "/") {
		pattern = Chop(pattern)
	}
	if strings.HasSuffix(path, "/") {
		path = Chop(path)
	}
	regStr := "^"

	// Do some syntax checking on the pattern.
	// filepath's Match() has some really weird rules that are inconsistent
	// so instead of trying to dup their logic, just call Match() for its
	// error state and if there is an error in the pattern return it.
	// If this becomes an issue we can remove this since its really only
	// needed in the error (syntax) case - which isn't really critical.
	if _, err := filepath.Match(pattern, path); err != nil {
		return false, err
	}

	// Go through the pattern and convert it to a regexp.
	// We use a scanner so we can support utf-8 chars.
	var scan scanner.Scanner
	scan.Init(strings.NewReader(pattern))

	sl := string(os.PathSeparator)
	escSL := sl
	if sl == `\` {
		escSL += `\`
	}

	for scan.Peek() != scanner.EOF {
		ch := scan.Next()

		if ch == '*' {
			if scan.Peek() == '*' {
				// is some flavor of "**"
				scan.Next()

				if scan.Peek() == scanner.EOF {
					// is "**EOF" - to align with .gitignore just accept all
					regStr += ".*"
				} else {
					// is "**"
					regStr += "((.*" + escSL + ")|([^" + escSL + "]*))"
				}

				// Treat **/ as ** so eat the "/"
				if string(scan.Peek()) == sl {
					scan.Next()
				}
			} else {
				// is "*" so map it to anything but "/"
				regStr += "[^" + escSL + "]*"
			}
		} else if ch == '?' {
			// "?" is any char except "/"
			regStr += "[^" + escSL + "]"
		} else if strings.Index(".$", string(ch)) != -1 {
			// Escape some regexp special chars that have no meaning
			// in golang's filepath.Match
			regStr += `\` + string(ch)
		} else if ch == '\\' {
			// escape next char. Note that a trailing \ in the pattern
			// will be left alone (but need to escape it)
			if sl == `\` {
				// On windows map "\" to "\\", meaning an escaped backslash,
				// and then just continue because filepath.Match on
				// Windows doesn't allow escaping at all
				regStr += escSL
				continue
			}
			if scan.Peek() != scanner.EOF {
				regStr += `\` + string(scan.Next())
			} else {
				regStr += `\`
			}
		} else {
			regStr += string(ch)
		}
	}

	regStr += "$"

	res, err := regexp.MatchString(regStr, path)

	// Map regexp's error to filepath's so no one knows we're not using filepath
	if err != nil {
		err = filepath.ErrBadPattern
	}

	return res, err
}
