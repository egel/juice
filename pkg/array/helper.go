package array

import (
	"fmt"
)

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

func Print(arr []string) {
	for i, val := range arr {
		fmt.Printf("%d.\t%s\n", i, val)
	}
}
