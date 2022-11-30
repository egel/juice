package multistring

import (
	"bufio"
	"regexp"
	"strings"
)

func RemoveProjectPackageNameLine(input string, packageName string) string {
	regProjectName := regexp.MustCompile("(?m)^" + packageName + ".*$")
	return regProjectName.ReplaceAllString(input, "")
}

func RemoveDedupedPackages(input string) string {
	regDeduped := regexp.MustCompile("(?m)[\r\n]+^.*deduped.*$")
	return regDeduped.ReplaceAllString(input, "")
}

func RemoveEmptyLines(input string) string {
	return regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(input), "\n")
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func MultilinestringToArray(multiline string) []string {
	var arr []string
	scanner := bufio.NewScanner(strings.NewReader(multiline))
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	return arr
}
