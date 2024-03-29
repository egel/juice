package npm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	url2 "net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	NodeModulesDirPath      = "./node_modules"
	PackageJsonFilePath     = "./package.json"
	PackageLockJsonFilePath = "./package-lock.json"
)

// one per single command, so it's parallel save
var execCommand_NpmList = exec.Command

type License struct {
	Name          string `json:"name"`    // source: npm field
	Version       string `json:"version"` // source: npm field
	LicenseType   string `json:"license"` // source: npm field
	LicenseText   string // source: node_modules
	LicenseUrl    string // source: experimental prediction of URL
	Homepage      string `json:"homepage"`       // source: npm field
	RepositoryUrl string `json:"repository.url"` // source: npm field
	NpmPackageUrl string // source: npm URL
	Error         string // source: used only for reporting errors while fetching
}

func RemoveAllNpmTreeCharacters(input string) string {
	str := strings.ReplaceAll(input, "│", "")
	str = strings.ReplaceAll(str, "│", "")
	str = strings.ReplaceAll(str, "├", "")
	str = strings.ReplaceAll(str, "─", "")
	str = strings.ReplaceAll(str, "┬", "")
	str = strings.ReplaceAll(str, "└", "")
	str = strings.ReplaceAll(str, " ", "")
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

func IsNodeModuleExist() bool {
	var path = NodeModulesDirPath
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Printf("directory not exist: '%s'", path)
		} else {
			log.Printf("general path error: '%v'", err)
		}
		return false
	}
	return true
}

func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetListProductionPackagesFromPackageLock() string {
	cmd := execCommand_NpmList("npm", "ls", "--production", "--all", "--package-lock-only")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println("can not get the list of production packages from package-lock.json")
		log.Fatal(err)
	}
	return string(stdout[:])
}

func PrintNumberOfPackages(arr []string) {
	fmt.Printf("Total number of production packages: %d\n", len(arr))
}

func FetchPackagesLicences(packageList []string) []License {
	maxConcurentRoutines := 80 // max amount of concurent routines
	packages := make(chan string, maxConcurentRoutines)
	licenses := make(chan License)
	var results []License
	count := len(packageList)

	// create array of channels
	for i := 0; i < cap(packages); i++ {
		go fetchPackageDetailsWorker(packages, licenses)
	}

	// spawn requesting packages
	go func() {
		for i := 0; i < count; i++ {
			packages <- packageList[i]
		}
	}()

	// save results to array
	for i := 1; i <= count; i++ {
		license := <-licenses
		fmt.Printf("\rProgress... %d/%d", i, count)
		results = append(results, license)
	}

	close(packages)
	close(licenses)

	return results
}

// this is worker function
func fetchPackageDetailsWorker(packageNames chan string, licenseResults chan License) {
	for p := range packageNames {
		var license License
		//fmt.Printf("\nfetching data for %s (%d)", p, len(packageNames))

		cmd := exec.Command("npm", "info", p, "--json", "name", "version", "license", "homepage", "repository.url")
		stdout, err := cmd.Output()
		if err != nil {
			msg := fmt.Sprintf("\n'%s': cmd fatal: %v", p, err)
			license.Error = msg
			licenseResults <- license
			break
		}
		err = json.Unmarshal(stdout, &license)
		if err != nil {
			msg := fmt.Sprintf("\n'%s': can not unmarshal: %v", p, err)
			license.Error = msg
			licenseResults <- license
			break
		}

		// Adding link to npm package
		license.NpmPackageUrl = "https://www.npmjs.com/package/" + license.Name + "/v/" + license.Version

		// Clean repository URL
		license.RepositoryUrl = ungitRepositoryUrl(license.RepositoryUrl)

		// Try getting LICENSE file from node_modules directory
		licenseFilePath, licenseFileName, err := checkExistenceOfLicenceFile("./node_modules/" + license.Name)
		if err != nil {
			msg := fmt.Sprintln(err)
			license.Error = msg
		} else {
			// Add license file (if exist)
			contents, err := os.ReadFile(licenseFilePath)
			if err != nil {
				msg := fmt.Sprintf("\n'%s': file reading error: %v", p, err)
				license.Error = msg
			}
			license.LicenseText = string(contents)

			// Add experimental link to license
			license.LicenseUrl = getLicenceFileUrl(license, licenseFileName)
		}

		licenseResults <- license
	}
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
	repoUrlClean := ungitRepositoryUrl(license.RepositoryUrl)
	repoUrlParsed, err := url2.Parse(repoUrlClean)
	if err != nil {
		fmt.Println("can not parse url", repoUrlClean, err)
	}
	var repoUrl string
	domain := repoUrlParsed.Host
	path := repoUrlParsed.Path

	switch domain {
	case "github.com":
		// The solution below tries to predict the most popular way by appending "v" to its numeric version (exposed by GitHub releases)
		// Although each repository may have different way of marking/tagging its versions, therefore some links may not work as expected.
		repoUrl = "https://raw.githubusercontent.com" + path + "/v" + license.Version + "/" + licenseFileName
	case "gitlab.com":
		repoUrl = "gitlab"
	default:
		repoUrl = "not known"
	}
	return repoUrl
}

func ungitRepositoryUrl(input string) string {
	result := input
	// prefix
	noGitPrefix := regexp.MustCompile(`(?m)^git\\+`)
	result = noGitPrefix.ReplaceAllString(input, "")
	// postfix
	noGitPostfix := regexp.MustCompile(`(?m)\\.git$`)
	result = noGitPostfix.ReplaceAllString(result, "")
	return result
}
