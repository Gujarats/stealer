package stealer

import (
	"fmt"
	"testing"
)

// Test if we have the all the variables name in the map key
func TestGetVariablesValue(t *testing.T) {
	testObjects := []struct {
		Datas        []byte
		ExpectedKeys []string
	}{
		// Test 0
		{
			Datas: []byte(`class php extend something {
				protected static $persons = array('andy','budy','yudi');
			}`),
			ExpectedKeys: []string{"persons"},
		},

		// Test 1
		{
			Datas: []byte(`

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
			`),
			ExpectedKeys: []string{"persons", "persons_static", "animals", "animals_static", "last_name", "last_name_static"},
		},
	}

	for indexTest, testObject := range testObjects {
		result := getVariablesValues(testObject.Datas)
		for key, _ := range result {
			exist := isValueExistArray(key, testObject.ExpectedKeys)
			if !exist {
				t.Errorf("index = %v, key = %v is not exist in expected keys =%+v\n", indexTest, key, testObject.ExpectedKeys)
			}
		}
	}
}

//Test if we have collected all correct values form each of variable
func TestFindData(t *testing.T) {
	testObjects := []struct {
		access   string
		data     []byte
		expected []string
	}{
		// Test 0
		{
			access: "private",
			data: []byte(`

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
			`),
			expected: []string{"andy", "clara", "john", "andys", "claras", "johns", "Abrahams", "Santanas", "Wijayas"},
		},

		// Test 1
		{
			access: "protected",
			data: []byte(`

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
			`),
			expected: []string{"Abraham", "Santana", "Wijaya"},
		},

		// Test 2
		{
			access: "public",
			data: []byte(`

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
			`),
			expected: []string{"lion", "wolf", "tiger", "lions", "wolfs", "tigers"},
		},

		//Test 3
		{
			access: "public",
			data: []byte(`

				public function __construct(){

				}

				private function someFunctionHere(){
					return 1;
				}

				public function publicFunctionHere(){
					return 1;
				}

    			protected static $tollFreeAreaCodes = array(
        			800, 844, 855, 866, 877, 888
    			);

				public $myFloatNumber = array(1.233, 123.1283, 11.22);

			`),
			expected: []string{"1.233", "123.1283", "11.22"},
		},

		//Test 4
		{
			access: "protected",
			data: []byte(`

				public function __construct(){

				}

				private function someFunctionHere(){
					return 1;
				}

				public function publicFunctionHere(){
					return 1;
				}

    			protected static $tollFreeAreaCodes = array(
        			800, 844, 855, 866, 877, 888
    			);

				public $myFloatNumber = array(1.233, 123.1283, 11.22);

			`),
			expected: []string{"800", "844", "855", "866", "877", "888"},
		},
	}

	for indexTest, testObject := range testObjects {
		actual := findData(testObject.access, testObject.data)
		for _, value := range testObject.expected {
			exist := isValueExist(value, actual)
			if !exist {
				fmt.Printf("index = %v, value = %v is not exist in actual = %+v\n", indexTest, value, actual)
			}
		}
	}
}

func TestGetValues(t *testing.T) {
	testObjects := []struct {
		index        int
		idxSemiColon int
		data         []byte
		expected     []string
	}{
		//Test 0
		{
			index:    0,
			data:     []byte(`private static $myVariable = array("data1","data2","data3");`),
			expected: []string{"data1", "data2", "data3"},
		},

		//Test 1
		{
			index:    0,
			data:     []byte(`private static $myVariable = array('one','two','three','four');`),
			expected: []string{"one", "two", "three", "four"},
		},

		//Test 2
		{
			index:    20,
			data:     []byte(`private static $myVariable = array('one','two','three','four');`),
			expected: []string{"one", "two", "three", "four"},
		},

		// Test 3
		{
			index:    20,
			data:     []byte(`private static $myVariable = array('one\'t','two','three\'okay','four\'hehe');`),
			expected: []string{"one't", "two", "three'okay", "four'hehe"},
		},

		// Test 4
		{
			index:    20,
			data:     []byte(`private static $myVariable = array("one\"t","two","three\"okay","four\"hehe","five\"");`),
			expected: []string{"one\\\"t", "two", "three\\\"okay", "four\\\"hehe", "five\\\""},
		},
	}

	for indexTest, testObject := range testObjects {
		_, actualValues := getValues(testObject.index, len(testObject.data)-1, testObject.data)
		ok := isSliceEqual(testObject.expected, actualValues)
		if !ok {
			t.Errorf("index test = %v, expected = %+v\n, actual = %+v\n", indexTest, testObject.expected, actualValues)
		}
	}
}

func TestAddValue(t *testing.T) {
	testObjects := []struct {
		store        []string
		data         []byte
		currentIndex int
		lastIndex    int
		expected     []string
	}{
		// test 0
		{
			store:        []string{"a", "b", "c"},
			data:         []byte(`hello world bother`),
			currentIndex: 0,
			lastIndex:    5,
			expected:     []string{"a", "b", "c", "hello"},
		},
		// test 1
		{
			store:        []string{"a", "b", "c"},
			data:         []byte(`f`),
			currentIndex: 0,
			lastIndex:    1,
			expected:     []string{"a", "b", "c", "f"},
		},
		// test 2
		{
			store:        []string{"a", "b", "c"},
			data:         []byte(`f      `),
			currentIndex: 0,
			lastIndex:    5,
			expected:     []string{"a", "b", "c", "f"},
		},
		// test 3
		{
			store:        []string{"a", "b", "c"},
			data:         []byte(`12      `),
			currentIndex: 0,
			lastIndex:    5,
			expected:     []string{"a", "b", "c", "12"},
		},
	}

	for index, testObject := range testObjects {
		addValueStore(&testObject.store, testObject.data, testObject.currentIndex, testObject.lastIndex)
		if !isSliceEqual(testObject.store, testObject.expected) {
			t.Errorf("at index %+v expected = %+v, result = %+v\n", index, testObject.expected, testObject.store)
		}
	}
}

func TestFindChar(t *testing.T) {
	testObjects := []struct {
		data         []byte
		char         []byte
		currentIndex int
		expected     int
	}{
		// test 0
		{
			data:         []byte(`hello world`),
			char:         []byte(`o`),
			currentIndex: 0,
			expected:     4,
		},
		// test 1
		{
			data:         []byte(`hello world`),
			char:         []byte(`d`),
			currentIndex: 5,
			expected:     10,
		},
		// test 2
		{
			data:         []byte(`hello world`),
			char:         []byte(`d`),
			currentIndex: 10,
			expected:     10,
		},
		// test 3
		{
			data:         []byte(`hello world`),
			char:         []byte(`d`),
			currentIndex: 100,
			expected:     -1,
		},
		// test 4
		{
			data:         []byte(`hello world`),
			char:         []byte(`z`),
			currentIndex: 0,
			expected:     -1,
		},
	}

	for index, testObject := range testObjects {
		result := findChar(testObject.data, testObject.char, testObject.currentIndex)
		if result != testObject.expected {
			t.Errorf("at index = %+v, expected = %+v, result = %+v\n", index, testObject.expected, result)
		}
	}
}

func isSliceEqual(a, b []string) bool {
	if a == nil && b != nil {
		return false
	}

	if a != nil && b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// checking if value a is exit in b map
func isValueExist(a string, b map[string][]string) bool {
	exist := false
mapLoop:
	for _, values := range b {
	valuesLoop:
		for _, value := range values {
			if value == a {
				exist = true
				break valuesLoop
				break mapLoop
			}
		}
	}
	return exist
}

// checking if value a is exit in b slice
func isValueExistArray(a string, b []string) bool {
	exist := false
	for _, value := range b {
		if value == a {
			exist = true
			break
		}
	}
	return exist
}
