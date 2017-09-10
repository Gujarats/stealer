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

You will get them in `map[string][]string` the key is variable name and the array is all the values.

## Usage

```go
import(
    "github.com/Gujarats/stealer"
)

//You wil get the result all the variable and its values from php in new path file written in go.
func main(){
    steal,err := stealer.Steal("path-php-file")
    if err != nil {
             log.Println(err)
    }
    steal.Save("test.go","package-name")
}

```

