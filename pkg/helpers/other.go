package helpers

import (
	"net/url"
	"strconv"
	"strings"
)

func ContainsNonEmptyString(s []string) bool {
	for _, v := range s {
		if v != "" {
			return true
		}
	}
	return false
}

func CountNewLines(s string) int {
	return countRune(s, '\n')
}

func countRune(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}

func Int64ArrayToJoinedString(a []int64) string {
	sArr := []string{}
	for _, v := range a {
		sArr = append(sArr, strconv.FormatInt(v, 10))
	}
	return strings.Join(sArr, ",")
}

/*
This function only exists for my own usage, I don't want to open some links in
the default browser.
*/
func OpenInChrome(url *url.URL) {
	CmdExec("cmd", "/C", "start", "chrome.exe", url.String())
}
