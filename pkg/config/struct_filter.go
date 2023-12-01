package config

import "fmt"
import "strings"
import "regexp"

// FilterConfig is a struct describing a single filter configuration for droping the access log line matched the filter
type FilterConfig struct {
	SourceValue    string `hcl:"from" yaml:"from"`
	Match          string `hcl:"match" yaml:"match"`
	Split          int    `hcl:"split" yaml:"split"`
	Separator      string `hcl:"separator" yaml:"seperator"`
	CompiledRegexp *regexp.Regexp
}

// Compile compiles expressions and lookup tables for efficient later use
func (fc *FilterConfig) Compile() error {
	if fc.Match != "" {
		r, err := regexp.Compile(fc.Match)
		if err != nil {
			return fmt.Errorf("could not compile regexp '%s': %s", fc.Match, err.Error())
		}

		fc.CompiledRegexp = r
	}

	return nil
}

func (fc FilterConfig) Filter(line string) bool {
	if fc.Split > 0 {
		separator := fc.Separator
		if separator == "" {
			separator = " "
		}

		values := strings.Split(line, separator)

		if len(values) >= fc.Split {
			line = values[fc.Split-1]
		} else {
			line = ""
		}
	}

	if fc.CompiledRegexp != nil {
		if fc.CompiledRegexp.MatchString(line) {
			return true
		}
	}

	return false

}
