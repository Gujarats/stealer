package main

import "fmt"

type testStruct struct {
	Data []interface{}
}

func main() {
	h := testStruct{
		Data: []interface{}{"check", []string{"h", "i"}},
	}

	testData := []interface{}{"check", []string{"h", "i"}}

	if h.Data[0] == "check" {
		fmt.Println("berhasil")
		return
	}

	size := 64
	data := make([]string, size)
	data[0] = "something"
	fmt.Println(len(data))
	fmt.Println(testData[0])
	fmt.Println(h.Data[1])
}
