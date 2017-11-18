package main

import (
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/Gujarats/stealer"
)

// convert all files in the given location
func fileConverter(location, fileSavePath string) {
	allFiles := getFiles(location)
	convertFilesToGo(allFiles, fileSavePath)
}

func convertFilesToGo(filesPath []string, filePathSaved string) {
	for _, filePath := range filesPath {
		convertFileToGO(filePath, filePathSaved)
	}
}

func convertFileToGO(filePath, fileSavePath string) {
	err, steal := stealer.Steal(filePath)
	if err != nil {
		log.Printf("couldnot steal data from given file = %+s and got error = %+v\n", filePath, err)
		return
	}

	fileName := convertFileName(filePath)

	// Save it to your path with the package name.
	packageName := getPackageName(fileSavePath)
	err = steal.Save(fileSavePath+fileName, packageName)
	if err != nil {
		log.Printf("couldnot save stole data got error = %+s\n", err)
		return
	}
}

// convert given path that has file name to Go's file extension
// return empty string if the file is not supported eg : support PHP
func convertFileName(filePath string) string {
	filePathSplitted := strings.Split(filePath, "/")
	if len(filePathSplitted) > 0 {
		fileName := filePathSplitted[len(filePathSplitted)-1]
		finalFileName := convertExtension(fileName)
		return finalFileName
	}

	return ""
}

func convertExtension(fileName string) string {
	fileNameSplitted := strings.Split(fileName, ".")

	if len(fileNameSplitted) > 1 && isExtentionSuported(fileNameSplitted[1]) {
		result := fileNameSplitted[0] + ".go"
		return result
	}

	return ""
}

func isExtentionSuported(extention string) bool {
	var supporExtensions = []string{"php"}
	var found bool
	for _, supporExtension := range supporExtensions {
		if supporExtension == extention {
			found = true
		}
	}

	return found
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
