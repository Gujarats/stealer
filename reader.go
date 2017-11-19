package stealer

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/Gujarats/logger"
)

// Read file and return its content.
func ReadFile(filepath string) (error, map[string][]string) {
	var result map[string][]string

	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err, nil
	}

	fileInfo, err := file.Stat()
	size := fileInfo.Size()
	datas := make([]byte, size)

	n, err := file.Read(datas)
	if err != nil {
		return err, nil
	}

	if n == 0 {
		err = errors.New("Empty selected file")
		return err, nil
	}

	result = getVariablesValues(datas)

	return nil, result
}

// Getting all the variables values that has prefix (protected,private,public)
// Return map with the key for variable name and []string for all the values from the variable
func getVariablesValues(datas []byte) map[string][]string {
	result := make(map[string][]string)
	accessKeys := []string{"private", "protected", "public"}
	for _, access := range accessKeys {
		datas := findData(access, datas)

		if len(datas) > 0 {
			// adding map datas to result
			for key, value := range datas {
				result[key] = value
			}
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

					if len(values) > 0 {
						// found values insert to our map
						result[varName] = values
						i = index
					} else {
						// we did'nt found any values from the variable, skip it
						i = varEndIndex
					}

				}
			}

		} else {
			break
		}

	}
	return result
}

// Getting all values number from a variable.
// the number could be any data type like eg: int,float
// call this function if we got null values from getValuesString()
// Note i is current index in integer
func getValuesNumber(i, idxSemiColon int, data []byte) (int, []string) {
	var valuesNumber []string
	comma := []byte(`,`)
	oP := []byte(`(`) // oP -> open Parenthesis
	cP := []byte(`)`) // cP ->close Parenthesis
	var cPindex int
	var oPindex int
	var commaIndex int

	// finding the open parenthesis `(` character and assign tu current index = i
	logger.Debug("currentInde = ", i)
	logger.Debug("len data = ", len(data))
	logger.Debug("data = ", string(data))
	logger.Debug("data current Index = ", string(data[i:i+20]))
	oPindex = findChar(data, oP, i)

	// check oPindex
	if oPindex > -1 {
		i = oPindex + 1
	}

	if oPindex >= idxSemiColon || oPindex == -1 {
		return i, valuesNumber
	}

	logger.Debug("opIndex = ", oPindex)
	logger.Debug("opIndex data = ", string(data[oPindex]))

	for {
		commaIndex = findChar(data, comma, i)
		logger.Debug("commaIndex = ", commaIndex)
		if commaIndex >= idxSemiColon {
			break
		}
		if commaIndex > -1 {
			// found separator get the value
			// set and check separatorIndex
			addValueStore(&valuesNumber, data, i, commaIndex)
			i = commaIndex + 1

			if i >= idxSemiColon {
				break
			}
		} else {
			// find close Parenthesis
			cPindex = findChar(data, cP, i)
			if cPindex >= idxSemiColon || cPindex == -1 {
				break
			}

			addValueStore(&valuesNumber, data, i, cPindex)
			logger.Debug("LAST valuesNumber= ", valuesNumber)
			i = cPindex + 1

		}

	}

	return i, valuesNumber
}

// Adding value to store from given data within the currentIndex and lastIndex
// remove the whitespaces if exist
// and return the currentIndex which modified from this function
func addValueStore(store *[]string, data []byte, currentIndex int, lastIndex int) {
	if currentIndex < lastIndex && lastIndex <= len(data) {
		value := data[currentIndex:lastIndex]
		if string(value) != "" {
			valueString := strings.TrimSpace(string(value))
			*store = append(*store, valueString)
		}
	}
}

// find the char index inside data from the given of currentIndex
// This currentIndex is to avoid getting previous charIndex
// return currentIndex and charIndex in sequence
func findChar(data []byte, char []byte, currentIndex int) int {
	charIndex := -1
	if currentIndex < len(data) {
		charIndex = bytes.Index(data[currentIndex:], char)
		if charIndex > -1 {
			//found char  inside data
			charIndex = currentIndex + charIndex
		}
	}

	return charIndex
}

// Getting the values from data with index i to idxSemiColon
// In php you can use single "'" and double "'".
// This function will get the values concurrently for the string values
// if string values not exist then it will try to get the numbers
func getValues(i, idxSemiColon int, data []byte) (int, []string) {
	// type and variable to hold the result
	type storage struct {
		Index  int
		Values []string
	}
	var finalResult storage

	// variable for reading the string concurrently
	var wg sync.WaitGroup
	quotes := [2]interface{}{[]byte(`"`), []byte(`'`)}
	values := make(chan storage, len(quotes))

	// checking for single quote to remove escapse
	escapse := []byte(`\`)
	singleQuote := []byte(`'`)
	removeEscape := false

	for _, quote := range quotes {
		wg.Add(1)
		go func(i, idxSemiColon int, quote []byte, data []byte) {
			defer wg.Done()
			var result []string
			for {
				removeEscape = false
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

				// finding the esacpe "\".
				if data[secondSep-1] == escapse[0] {
					// if we escaping single quote
					if data[secondSep] == singleQuote[0] {
						// mark to remove escapse
						removeEscape = true
					}

					secondSep = bytes.Index(data[i+1:], quote)
					i = i + secondSep + 1
					secondSep = i
				}

				value := data[firstSep:secondSep]

				if removeEscape {
					var newValue []byte
					for _, char := range value {
						if char != escapse[0] {
							newValue = append(newValue, char)
						}
					}
					value = newValue
				}
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

	var foundStringValue bool
	for value := range values {
		if len(value.Values) > 0 {
			// found values
			foundStringValue = true
			finalResult = value
		}
	}

	if !foundStringValue {
		// not found values using quotes eg : "\'" and "\""
		// it could be number
		index, valuesNumber := getValuesNumber(i, idxSemiColon, data)
		if len(valuesNumber) > 0 {
			finalResult.Index = index
			finalResult.Values = valuesNumber
		}

	}

	return finalResult.Index, finalResult.Values
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
