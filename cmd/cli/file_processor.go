package main

import (
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/Gujarats/stealer"
)

// convert all files in the given location
func fileConverter(location, filePathSaved string) {
	allFiles := getFiles(location)
	convertFilesToGo(allFiles, filePathSaved)
}

func convertFilesToGo(filesPath []string, filePathSaved string) {
	for _, filePath := range filesPath {
		convertFileToGO(filePath, filePathSaved)
	}
}

func convertFileToGO(filePath, filePathSaved string) {
	err, steal := stealer.Steal(filePath)
	if err != nil {
		log.Printf("couldnot steal data from given file = %+s and got error = %+v\n", filePath, err)
		return
	}

	// Save it to your path with the package name.
	packageName := getPackageName(filePathSaved)
	err = steal.Save(filePathSaved, packageName)
	if err != nil {
		log.Println("couldnot save stole data got error = %+v\n", err)
		return
	}
}

// TODO : move this function inside steal.Save() ?
// get the package name from given filepathLocation
// eg : my/path/location/file.go the package name will be "location"
func getPackageName(filePath string) string {
	splitPath := strings.Split(filePath, "/")
	return splitPath[len(splitPath)-2]
}

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

// check if file is php or not
func checkFileExtension(filePath string) error {
	return nil
}
