package stealer

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

// Write all the variables and arrays into code to save them for the later usage
func Write(path, packageName string, data map[string][]string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	//close file after finish writing
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	// writing package name
	file.WriteString(packageName + "\n")

	// writing variable
	for varName, values := range data {
		varible := variableFormat(varName, values)
		if _, err := file.WriteString(varible); err != nil {
			log.Println(err)
		}
	}

	return nil
}

// the varibale format should be varName := []interface{values}
func variableFormat(varName string, values []string) string {
	// get values data type
	dataType := dataTypeFormat(values)

	var variable string
	variable = varName
	variable += " := []" + dataType + "{"

	if dataType == "bool" || dataType == "int" || dataType == "float64 " {
		for index, value := range values {
			if index != len(values)-1 {
				variable += value + ","
			} else {
				variable += value
			}
		}
		variable += "}"
	} else {
		// found string data type
		for index, value := range values {
			if index != len(values)-1 {
				variable += "\"" + value + "\","
			} else {
				variable += "\"" + value + "\""
			}
		}
		variable += "}"
	}

	return variable
}

// Assuming that values []string are all the same data type
func dataTypeFormat(values []string) string {
	// check wheter values are boolean or not
	var result interface{}
	var err error
	foundType := make(map[bool]bool)
	for _, value := range values {
		result, err = strconv.ParseBool(value)
		if err == nil {
			foundType[true] = true
		} else {
			foundType[false] = true
		}
	}

	if !foundType[false] {
		return fmt.Sprintf("%s", reflect.TypeOf(result))
	}

	result, err = strconv.ParseInt(values[0], 10, 64)
	if err == nil {
		return "int"
	}

	_, err = strconv.ParseFloat(values[0], 64)
	if err == nil {
		return "float64"
	}

	return "string"
}
