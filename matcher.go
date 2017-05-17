package main

import (
	"regexp"
	"strings"
)

type matcherFunc func([]string, string, []string) []string

func (mf matcherFunc) Match(args []string, pat string, list []string) []string {
	return mf(args, pat, list)
}

func init() {
	registerMatcher("string", matcherFunc(matchString))
	registerMatcher("substring", matcherFunc(matchSubstring))
	registerMatcher("token", matcherFunc(matchToken))
	registerMatcher("subtoken", matcherFunc(matchSubtoken))
	registerMatcher("alias_regexp", matcherFunc(matchAliasRegexp))
}

// string match, pattern is exactly the same as target.
func matchString(args []string, pat string, list []string) []string {
	var result []string
	for _, v := range list {
		if pat == v {
			result = append(result, v)
		}
	}
	return result
}

// pattern is a substring of target.
func matchSubstring(args []string, pat string, list []string) []string {
	var result []string
	for _, v := range list {
		if strings.Contains(v, pat) {
			result = append(result, v)
		}
	}
	return result
}

// match token separated by option
func matchToken(args []string, pat string, list []string) []string {
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
func matchSubtoken(args []string, pat string, list []string) []string {
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

func matchAliasRegexp(args []string, pat string, list []string) []string {
	if len(args) < 2 {
		return nil
	}

	hostPattern := args[0]
	template := args[1]

	re, err := regexp.Compile(hostPattern)
	if err != nil {
		return nil
	}
	match := re.FindStringSubmatchIndex(pat)
	if match == nil {
		return nil
	}

	exp := re.ExpandString(nil, template, pat, match)

	return []string{string(exp)}
}
