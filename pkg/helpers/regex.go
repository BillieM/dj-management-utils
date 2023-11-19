package helpers

import "regexp"

func regexReplace(s string, regexStr string, replacementStr string) string {
	re := regexp.MustCompile(regexStr)
	s = re.ReplaceAllString(s, replacementStr)
	return s
}

func regexContains(s string, regexStr string) bool {
	re := regexp.MustCompile(regexStr)
	return re.MatchString(s)
}
