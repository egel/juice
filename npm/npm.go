package npm

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type License struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	License       string `json:"license"`
	Homepage      string `json:"homepage"`
	RepositoryUrl string `json:"repository.url"`
	LicenseText   string
}

func RemoveAllNpmTreeCharacters(input string) string {
	str := strings.Replace(input, "│", "", -1)
	str = strings.Replace(str, "│", "", -1)
	str = strings.Replace(str, "├", "", -1)
	str = strings.Replace(str, "─", "", -1)
	str = strings.Replace(str, "┬", "", -1)
	str = strings.Replace(str, "└", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

func GetPackageDetails(ch chan<- License, packageName string) {
	fmt.Println("starting " + packageName)
	var license License
	cmd := exec.Command("npm", "info", packageName, "--json", "name", "version", "license", "homepage", "repository.url")

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd fatal: %v", err)
	}

	err = json.Unmarshal(stdout, &license)
	if err != nil {
		fmt.Println("can not unmarshal!", err)
	}

	// TODO: get licence text

	ch <- license
	close(ch)
}
