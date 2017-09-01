package stealer

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

// Getting all the variables values that has prefix (protected,private,public)
// Return map with the key for variable name and []string for all the values from the variable
func GetVariablesValues(datas []byte) map[string][]string {
	result := make(map[string][]string)
	accessKeys := []string{"private", "protected", "public"}
	for _, access := range accessKeys {
		datas := findData(access, datas)
		// adding map datas to result
		for key, value := range datas {
			result[key] = value
		}
	}

	return result
}

// Get all the values from a variable with specific access from the data.
// This acces can be private, protected or public.
func findData(access string, data []byte) map[string][]string {
	result := make(map[string][]string)
	function := []byte(`function`)
	lenAccess := len(access)

	for i := 0; i < len(data); {
		varIndex := bytes.Index(data[i:], []byte(access))
		if varIndex > -1 {
			i = i + varIndex + lenAccess
			// check if it is a variable or function
			indexEnter := bytes.Index(data[i:], []byte("\n"))
			funcIndex := bytes.Index(data[i:i+indexEnter], function)
			if funcIndex > -1 {
				// jump to specific index and continue search for variable
				i = i + indexEnter
				continue
			} else {
				// found a variable get variable name
				dollar := []byte(`$`)
				space := []byte(`=`)
				idxDollar := bytes.Index(data[i:], dollar)
				i = i + idxDollar
				idxSpace := bytes.Index(data[i:], space)
				varNameByte := data[i+1 : i+idxSpace]
				varName := removeSpace(string(varNameByte))

				//check if it is an array
				idxArray := bytes.Index(data[i:], []byte(`array`))
				if idxArray > -1 {
					i = i + idxArray
					// looop to get all the array
					varEndIndex := bytes.Index(data[i:], []byte(`;`))
					varEndIndex = i + varEndIndex
					index, values := getValues(i, varEndIndex, data)
					result[varName] = values
					i = index
				}
			}

		} else {
			break
		}

	}
	return result
}

// Getting the values from data with index i to idxSemiColon
// In php you can use single "'" and double "'".
// This function will get the values concurrently
func getValues(i, idxSemiColon int, data []byte) (int, []string) {
	var wg sync.WaitGroup
	quotes := [2]interface{}{[]byte(`"`), []byte(`'`)}
	type storage struct {
		Index  int
		Values []string
	}
	values := make(chan storage, len(quotes))

	for _, quote := range quotes {
		wg.Add(1)
		go func(i, idxSemiColon int, quote []byte, data []byte) {
			defer wg.Done()
			var result []string
			for {
				//getting all the values from array
				firstSep := bytes.Index(data[i:], quote)
				if firstSep == -1 {
					break
				}
				if firstSep+i >= idxSemiColon {
					i = idxSemiColon + 1
					break
				}

				i = i + firstSep + 1
				firstSep = i
				secondSep := bytes.Index(data[i:], quote)
				i = i + secondSep
				secondSep = i
				value := data[firstSep:secondSep]
				if string(value) == "" {
					break
				}
				result = append(result, string(value))
				i = i + 2
			}
			objectResult := storage{Index: i, Values: result}
			values <- objectResult
		}(i, idxSemiColon, quote.([]byte), data)
	}

	wg.Wait()

	close(values)
	var finalResult storage
	for value := range values {
		if value.Index != 0 && value.Values != nil {
			finalResult = value
		}
	}

	return finalResult.Index, finalResult.Values
}

// Read file and return its content.
func getData(filepath string) []byte {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	fileInfo, err := file.Stat()
	size := fileInfo.Size()
	datas := make([]byte, size)

	n, err := file.Read(datas)
	if err != nil {
		log.Fatalln(err)
	}
	if n == 0 {
		log.Fatalln(errors.New("Empty selected file"))
	}

	return datas
}

// Remove white space from string the fastest way
// See here https://stackoverflow.com/questions/32081808/strip-all-whitespace-from-a-string-in-golang
func removeSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
