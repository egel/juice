package npm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	url2 "net/url"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

const NodeModulesDirPath = "./node_modules"
const NpmPackageFilePath = "package.json"
const NpmPackageLockFilePath = "package-lock.json"

type License struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	LicenseType   string `json:"license"`
	LicenseText   string
	LicenseUrl    string
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
	if _, err := os.Stat(NodeModulesDirPath); !os.IsNotExist(err) {
		fmt.Println("Skipping installation and using existing node_modules (remove node_modules to trigger installation)")
		return
	} else {
		fmt.Println("(it may take up to few minutes...)")
	}

	cmd := exec.Command("npm", "ci")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Error while installing packages: ", err)
	}
}

func IsNpmPackagesFilesExistOrDie() {
	if _, err := exists(NpmPackageFilePath); err != nil {
		log.Fatalf("%s does not exist", NpmPackageFilePath)
	}
	if _, err := exists(NpmPackageLockFilePath); err != nil {
		log.Fatalf("%s does not exist", NpmPackageLockFilePath)
	}
}

func IsNodeModuleExistOrDie() {
	var path = NodeModulesDirPath
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("'%s' directory not exist", path)
		} else {
			log.Fatal("some other error", err)
		}
	}
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
	licenseFilePath, licenseFileName, err := checkExistenceOfLicenceFile("./node_modules/" + license.Name)
	if err != nil {
		fmt.Println(err)
	} else {
		// license file
		contents, err := os.ReadFile(licenseFilePath)
		if err != nil {
			fmt.Println("file reading error", err)
		}
		license.LicenseText = string(contents)

		// license link
		link := getLicenceFileUrl(license, licenseFileName)
		license.LicenseUrl = link

		// fix link for git repo
		repo := ungitRepoUrl(license.RepositoryUrl)
		license.RepositoryUrl = repo
	}

	ch <- license
	close(ch)
}

func checkExistenceOfLicenceFile(path string) (string, string, error) {
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
			return filePath, el, nil
		}
	}
	return "", "", errors.New("license file not found in path: " + path)
}

func getLicenceFileUrl(license License, licenseFileName string) string {
	repoUrlClean := ungitRepoUrl(license.RepositoryUrl)
	repoUrlParsed, err := url2.Parse(repoUrlClean)
	if err != nil {
		fmt.Println("can not parse url", repoUrlClean, err)
	}
	var repoUrl string
	domain := repoUrlParsed.Host
	path := repoUrlParsed.Path

	switch domain {
	case "github.com":
		repoUrl = "https://raw.githubusercontent.com" + path + "/" + license.Version + "/" + licenseFileName
	case "gitlab.com":
		repoUrl = "gitlab"
	default:
		repoUrl = "not known"
	}
	return repoUrl
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
