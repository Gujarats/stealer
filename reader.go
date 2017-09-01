package stealer

import (
	"bytes"
	"errors"
	"log"
	"os"
	"sync"
)

// checking if the variable has prefix (protected,private,public)
// and get the values from it
func GetVariablesValue(datas []byte) map[string]interface{} {
	result := make(map[string]interface{})
	//access := []string{"private", "protected", "public"}
	result["test"] = []interface{}{"hoho", []string{"hello", "world"}}

	return result
}

// get all the values from a variable with specific access from the data.
// this acces can be private, protected or public.
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
				varName := data[i+1 : i+idxSpace]

				//check if it is an array
				idxArray := bytes.Index(data[i:], []byte(`array`))
				if idxArray > -1 {
					i = i + idxArray
					// looop to get all the array
					varEndIndex := bytes.Index(data[i:], []byte(`;`))
					varEndIndex = i + varEndIndex
					index, values := getValues(i, varEndIndex, data)
					result[string(varName)] = values
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

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
