package helpers

import "regexp"

func RegexReplace(s string, regexStr string, replacementStr string) string {
	re := regexp.MustCompile(regexStr)
	s = re.ReplaceAllString(s, replacementStr)
	return s
}

func RegexContains(s string, regexStr string) bool {
	re := regexp.MustCompile(regexStr)
	return re.MatchString(s)
}

func RegexAllSubmatches(s string, regexStr string) [][]string {
	re := regexp.MustCompile(regexStr)
	return re.FindAllStringSubmatch(s, -1)
}
