# Stealer [![Build Status](https://secure.travis-ci.org/Gujarats/stealer.png)](http://travis-ci.org/Gujarats/stealer)
Get all the variable name and the values from php, for example if you have some php code file `mycode.php` like : 

```php
    public function __construct(){
    
    }
    
    private function someFunctionHere(){
        return 1;
    }
    
    public function publicFunctionHere(){
        return 1;
    }
    
    private $persons = array('andy','clara','john');
    private static $persons_static = array("andys","claras","johns");
    public $animals = array("lion","wolf","tiger");
    public static $animals_static = array("lions","wolfs","tigers");
    protected $last_name= array('Abraham','Santana','Wijaya');
    private static $last_name_static = array('Abrahams','Santanas','Wijayas');
```

This library will convert all those variable and its value to GO.

## CLI
You can directly use this using CLI in the `cmd/cli` directory, here is the step installation for Ubuntu : 

```shell
cd $GOPATH/src/github.com/Gujarats/stealer/cmd/cli/
go build
sudo mv cli /usr/local/bin/stealer
stealer
```
## Usage from the source
Here is some snippet code to get started

```go
package main

import (
	"log"

	"github.com/Gujarats/stealer"
)

func main() {
	err, steal := stealer.Steal("Person.php")
	if err != nil {
		log.Println(err)
	}

    // Save it to your path with the package name.
    // In this case the package name is main
	err = steal.Save("path/to/your/file/test.go", "main")
	if err != nil {
		log.Println(err)
	}
}

```
