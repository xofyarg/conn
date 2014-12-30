package main

import (
	"regexp"
	"strings"
)

type matcher struct {
	name  string
	match func(string, []string) []string
}

func (m *matcher) Match(pat string, list []string) []string {
	return m.match(pat, list)
}

func init() {
	registerMatcher("string", &matcher{"string", matchString})
	registerMatcher("substring", &matcher{"substring", matchSubstring})
	registerMatcher("token", &matcher{"token", matchToken})
	registerMatcher("subtoken", &matcher{"subtoken", matchSubtoken})
}

// string match, pattern is exactly the same as target.
func matchString(pat string, list []string) []string {
	var result []string
	for _, v := range list {
		if pat == v {
			result = append(result, v)
		}
	}
	return result
}

// pattern is a substring of target.
func matchSubstring(pat string, list []string) []string {
	var result []string
	for _, v := range list {
		if strings.Contains(v, pat) {
			result = append(result, v)
		}
	}
	return result
}

// match token separated by option
func matchToken(pat string, list []string) []string {
	keys := strings.Split(pat, options.Sep)
	var result []string
	re := regexp.MustCompile("[^.-]+")
	for _, v := range list {
		parts := re.FindAllString(v, -1)

		findMatch := true
		for _, key := range keys {
			match := false
			for _, part := range parts {
				if key == part {
					match = true
					break
				}
			}

			if !match {
				findMatch = false
				break
			}
		}

		if findMatch {
			result = append(result, v)
		}
	}

	return result
}

// match token by substring
func matchSubtoken(pat string, list []string) []string {
	keys := strings.Split(pat, options.Sep)
	var result []string
	re := regexp.MustCompile("[^.-]+")
	for _, v := range list {
		parts := re.FindAllString(v, -1)

		findMatch := true
		for _, key := range keys {
			match := false
			for _, part := range parts {
				if strings.Contains(part, key) {
					match = true
					break
				}
			}

			if !match {
				findMatch = false
				break
			}
		}

		if findMatch {
			result = append(result, v)
		}
	}

	return result
}
