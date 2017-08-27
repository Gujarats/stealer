package stealer

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
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
func findData(access string, data []byte) map[string][]interface{} {
	result := make(map[string][]interface{})
	function := []byte(`function`)
	lenAccess := len(access)
	//lenFunc := len(function)
firstLoop:
	for i := 0; i < len(data); {
		fmt.Println("masuk = ", i)
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
				LoopValue:
					for {

						//getting all the values from array
						sep := []byte(`'`)
						firstSep := bytes.Index(data[i:], sep)
						if firstSep == -1 {
							break firstLoop
						}
						if firstSep+i >= varEndIndex {
							i = varEndIndex + 1
							break LoopValue
						} else {

							i = i + firstSep + 1
							firstSep = i
							secondSep := bytes.Index(data[i:], sep)
							i = i + secondSep
							secondSep = i
							fmt.Println("first sep = ", firstSep)
							fmt.Println("second sep = ", secondSep)
							value := data[firstSep:secondSep]
							if string(value) != "" {
								fmt.Println("name - ", string(value))
								result[string(varName)] = append(result[string(varName)], string(value))
							}
							i = i + 2
						}
					}

				}
			}

		}

	}
	fmt.Printf("result exit loop = %+v\n", result)

	return result
}

// Getting the value from quoate
// In php you can use single "'" and double "'".
func getValue(i, idxSemiColon int, data []byte) (int, string, error) {
	var value string
	sep := []byte(`'`)
	firstSep := bytes.Index(data[i:idxSemiColon], sep)
	if firstSep == -1 {
		sep = []byte(`"`)
		firstSep = bytes.Index(data[i:idxSemiColon], sep)
		if firstSep == -1 {
			return i, value, errors.New("no data found")
		}
	}

	i = i + firstSep + 1
	firstSep = i

	if firstSep >= idxSemiColon {
		fmt.Println("varEndIndex = ", idxSemiColon)
		fmt.Println("i = ", i)
		i = i + idxSemiColon
		return i, value, errors.New("data out of limit")
	}

	return i, value, nil
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
