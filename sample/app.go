package main

import (
	"log"

	"github.com/gujarats/stealer"
)

func main() {
	err, steal := stealer.Steal("Person.php")
	if err != nil {
		log.Println(err)
	}
	err = steal.Save("someFolder/test/go/test.go", "main")
	if err != nil {
		log.Println(err)
	}
}
