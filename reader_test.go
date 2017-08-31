package stealer

import (
	"testing"
)

//func TestGetVariablesValue(t *testing.T) {
//	testObjects := []struct {
//		Datas    []byte
//		Expected []interface{}
//	}{
//		{
//			Datas: []byte(`class php extend something {
//				private static somevarible = "here";
//				protected static persons = array('andy','budy','yudi');
//			}`),
//			Expected: []interface{}{"here", []string{"andy", "budy", "yudi"}},
//		},
//	}
//
//	for _, testObject := range testObjects {
//		result := GetVariablesValue(testObject.Datas)
//		index := 0
//		for _, value := range result {
//			fmt.Println(testObject.Expected[1], index, value)
//
//			if testObject.Expected[index] != value {
//				t.Errorf("expected = %+v, actual = %+v\n", testObject.Expected[index], value)
//			}
//			index++
//		}
//	}
//}

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
			expected: []string{"andy", "clara", "john", "Abrahams", "Santanas", "Wijayas"},
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
	}

	for _, testObject := range testObjects {
		actual := findData(testObject.access, testObject.data)
		index := 0
		for _, actualValues := range actual {
			for _, actualValue := range actualValues {
				if testObject.expected[index] != actualValue {
					t.Errorf("actual = %v, expected = %v", actualValue, testObject.expected[index])
				}
				index++

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
	}

	for indexTest, testObject := range testObjects {
		_, actualValues := getValues(testObject.index, len(testObject.data)-1, testObject.data)
		ok := isSliceEqual(testObject.expected, actualValues)
		if !ok {
			t.Errorf("index test = %v, expected = %+v\n, actual = %+v\n", indexTest, testObject.expected, actualValues)
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
