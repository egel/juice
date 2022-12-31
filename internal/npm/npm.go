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
	"sort"
	"strings"
)

const NodeModulesDirPath = "./node_modules"
const PackageJsonFilePath = "./package.json"
const PackageLockJsonFilePath = "./package-lock.json"

type License struct {
	Name          string `json:"name"`    // npm field
	Version       string `json:"version"` // npm field
	LicenseType   string `json:"license"` // npm field
	LicenseText   string // taken from node_modules
	LicenseUrl    string // experimental prediction of URL
	Homepage      string `json:"homepage"`       // npm field
	RepositoryUrl string `json:"repository.url"` // npm field
	NpmPackageUrl string // npm URL
	Error         string // used only for reporting errors while fetching
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
	if _, err := exists(PackageJsonFilePath); err != nil {
		log.Fatalf("%s does not exist", PackageJsonFilePath)
	}
	if _, err := exists(PackageLockJsonFilePath); err != nil {
		log.Fatalf("%s does not exist", PackageLockJsonFilePath)
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
		log.Println("can not get the list of production packages from package-lock.json")
		log.Fatal(err)
	}
	return string(stdout[:])
}

func PrintNumberOfPackages(arr []string) {
	fmt.Printf("Total number of production packages: %d\n", len(arr))
}

func FetchPackagesLicences(packageList []string) []License {
	packages := make(chan string, 80) // max amount of concurent routines, currently fixed
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
		fmt.Printf("count: %d/%d\n", i, count)
		results = append(results, license)
	}

	close(packages)
	close(licenses)

	return results
}

func SortSlices(licenses []License) {
	sort.Slice(licenses, func(i, j int) bool {
		return licenses[i].Name < licenses[j].Name
	})
}

// this is worker function
func fetchPackageDetailsWorker(packageNames chan string, licenseResults chan License) {
	for p := range packageNames {
		var license License
		fmt.Printf("fetching data for %s (%d)\n", p, len(packageNames))

		cmd := exec.Command("npm", "info", p, "--json", "name", "version", "license", "homepage", "repository.url")
		stdout, err := cmd.Output()
		if err != nil {
			msg := fmt.Sprintf("'%s': cmd fatal: %v\n", p, err)
			license.Error = msg
			licenseResults <- license
			break
		}
		err = json.Unmarshal(stdout, &license)
		if err != nil {
			msg := fmt.Sprintf("'%s': can not unmarshal: %v\n", p, err)
			license.Error = msg
			licenseResults <- license
			break
		}

		// Adding link to npm package
		license.NpmPackageUrl = "https://www.npmjs.com/package/" + license.Name + "/v/" + license.Version

		// Clean repository URL
		license.RepositoryUrl = ungitRepositoryUrl(license.RepositoryUrl)

		// Get LICENSE
		licenseFilePath, licenseFileName, err := checkExistenceOfLicenceFile("./node_modules/" + license.Name)
		if err != nil {
			fmt.Println(err)
			license.Error = err.Error()
		} else {
			// Add license file (if exist)
			contents, err := os.ReadFile(licenseFilePath)
			if err != nil {
				msg := fmt.Sprintf("'%s': file reading error: %v\n", p, err)
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
