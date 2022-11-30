package npm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
)

type License struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	LicenseType   string `json:"license"`
	LicenseText   string
	Homepage      string `json:"homepage"`
	RepositoryUrl string `json:"repository.url"`
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

func InstallPackageLock() {
	fmt.Println("Installing packages from 'package-lock.json'")
	if _, err := os.Stat("./node_modules"); !os.IsNotExist(err) {
		fmt.Println("using cached version!")
		return
	}
	fmt.Println("it may take few minutes...")

	cmd := exec.Command("npm", "ci")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalln("can not install packages:", err)
	}
	fmt.Println(string(stdout))
}

func GetListProductionPackagesFromPackageLock() string {
	cmd := exec.Command("npm", "ls", "--production", "--all", "--package-lock-only")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(stdout[:])
}

func FetchPackagesLicences(packageList []string) []License {
	var licences []License

	count := len(packageList)
	//count := 5

	// create array of channels
	var myChannels []chan License
	for i := 0; i < count; i++ {
		ch := make(chan License)
		myChannels = append(myChannels, ch)
	}

	// create cases and spawn requests
	fmt.Println("Start spawning parallel requests and waiting for results...")
	cases := make([]reflect.SelectCase, len(myChannels))
	for k, ch := range myChannels {
		cases[k] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		go fetchPackageDetails(ch, packageList[k])
	}

	// remaining
	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining -= 1
			continue
		}

		licences = append(licences, value.Interface().(License))
	}

	return licences
}

func SortSlices(licenses []License) {
	sort.Slice(licenses, func(i, j int) bool {
		return licenses[i].Name < licenses[j].Name
	})
}

func fetchPackageDetails(ch chan<- License, packageName string) {
	fmt.Println("starting fetching data for " + packageName)
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

	// Get LICENSE
	filePath, err := checkExistenceOfLicenceFile("./node_modules/" + license.Name)
	if err != nil {
		fmt.Println(err)
	} else {
		contents, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("file reading error", err)
		}
		license.LicenseText = string(contents)
	}

	ch <- license
	close(ch)
}

func checkExistenceOfLicenceFile(path string) (string, error) {
	arr := []string{
		"LICENSE",
		"LICENSE.txt",
		"LICENSE.md",
		"COPYING",
		"COPYING.txt",
		"COPYING.md",
	}
	for _, el := range arr {
		filePath := path + "/" + el
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}
	}
	return "", errors.New("license file not found in path: " + path)
}
