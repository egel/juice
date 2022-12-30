package multistring

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

func RemoveProjectPackageNameLine(multiline string, packageName string) string {
	regProjectName := regexp.MustCompile("(?m)^" + packageName + ".*$")
	return regProjectName.ReplaceAllString(multiline, "")
}

func RemoveDedupedPackages(multiline string) string {
	regDeduped := regexp.MustCompile("(?m)[\r\n]+^.*deduped.*$")
	result := regDeduped.ReplaceAllString(multiline, "")
	fmt.Println("Deduped packages removed")
	return result
}

func RemoveEmptyLines(input string) string {
	return regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(input), "\n")
}

func MultilinestringToArray(multiline string) []string {
	var arr []string
	scanner := bufio.NewScanner(strings.NewReader(multiline))
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	return arr
}
