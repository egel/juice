package converter

import (
	"encoding/csv"
	"log"
	"os"

	"juice/pkg/npm"
)

func SaveDataToCSVFile(licenses []npm.License) {
	csvFile, err := os.Create("licences.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	writer := csv.NewWriter(csvFile)
	headers := []string{"Package Name", "Version", "Home URL", "Repository URL", "License Type", "License Text"}
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
		license.Homepage,
		license.RepositoryUrl,
		license.LicenseType,
		license.LicenseText,
	}
}
