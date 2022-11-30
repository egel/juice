package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"os/exec"
	"reflect"
	"regexp"
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

	//count := len(cleanPackages)
	count := 5

	// create array of channels
	var myChannels []chan License
	for i := 0; i < count; i++ {
		ch := make(chan License)
		myChannels = append(myChannels, ch)
	}

	// create cases and spawn requests
	cases := make([]reflect.SelectCase, len(myChannels))
	for k, ch := range myChannels {
		cases[k] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		go fetchPackageDetails(ch, packageList[k])
	}
	fmt.Println("Waiting for results...")

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
		license, _ := reflect.ValueOf(value).Interface().(License)
		licences = append(licences, license)
		fmt.Printf("Read from channel %#v and received %#v\n", myChannels[chosen], value)
	}

	return licences
}

func fetchPackageDetails(ch chan<- License, packageName string) {
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

	// Get LICENSE
	contents, err := os.ReadFile("./node_modules/" + license.Name + "/LICENSE")
	if err != nil {
		fmt.Println("File reading error", err)
	}
	license.LicenseText = string(contents)

	ch <- license
	close(ch)
}

func ungitRepoUrl(input string) string {
	result := input
	// prefix
	noGitPrefix := regexp.MustCompile("(?m)^git\\+")
	result = noGitPrefix.ReplaceAllString(input, "")
	// postfix
	noGitPostfix := regexp.MustCompile("(?m)\\.git$")
	result = noGitPostfix.ReplaceAllString(result, "")
	return result
}

// @deprecated
func getLicenceOld(license License) {
	// TODO: get licence text
	urlString := ungitRepoUrl(license.RepositoryUrl)
	fmt.Printf("urlString: %#v", urlString)

	url, err := url2.Parse(urlString)
	if err != nil {
		log.Fatal("fatal url parse", err)
	}
	domain := url.Host
	path := url.Path
	//fmt.Println("url.Host:", domain, path)

	var repoUrl string
	switch domain {
	case "github.com":
		fmt.Println("github")
		repoUrl = "https://raw.githubusercontent.com" + path + "/" + license.Version + "/LICENSE"
		fmt.Println(repoUrl)
	case "gitlab.com":
		fmt.Println("gitlab!")
	default:
		fmt.Println("INNE!", license)
	}

	resp, err := http.Get(repoUrl)
	if err != nil {
		fmt.Println("Error getting licence from:", repoUrl)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("can not read body:", err)
	}
	//Convert the body to type string
	sb := string(body)
	license.LicenseText = sb
}
