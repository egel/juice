package converter

import (
	"encoding/csv"
	"juice/internal/npm"
	"log"
	"os"
)

func SaveDataToCSVFile(licenses []npm.License) {
	csvFile, err := os.Create("licences.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	writer := csv.NewWriter(csvFile)
	headers := []string{
		"Package Name (npm)",
		"Version (npm)",
		"Link to NPM package (npm)", // FIXME
		"Homepage (npm)",
		"Repository URL (npm)",
		"License Type (npm)",
		"License Text (node_modules)",
		"License Url (experimental forecast, can be incorrect)",
		"Error",
	}
	if err := writer.Write(headers); err != nil {
		log.Fatal("can not save header for CSV file")
	}
	for _, license := range licenses {
		slice := structToSlice(license)
		if err := writer.Write(slice); err != nil {
			log.Fatalf("can not save package '%s' in CSV file", license.Name)
		}
	}
	writer.Flush()
	if err := csvFile.Close(); err != nil {
		log.Fatalf("file can not be closed")
	}
}

func structToSlice(license npm.License) []string {
	return []string{
		license.Name,
		license.Version,
		license.NpmPackageUrl,
		license.Homepage,
		license.RepositoryUrl,
		license.LicenseType,
		license.LicenseText,
		license.LicenseUrl,
		license.Error,
	}
}
