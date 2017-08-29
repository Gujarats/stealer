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
				private $static persons_static = array("andys","claras","johns");
				public $animals = array("lion","wolf","tiger");
				public $static animals_static = array("lions","wolfs","tigers");
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
				private $static persons_static = array("andys","claras","johns");
				public $animals = array("lion","wolf","tiger");
				public $static animals_static = array("lions","wolfs","tigers");
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
				private $static persons_static = array("andys","claras","johns");
				public $animals = array("lion","wolf","tiger");
				public $static animals_static = array("lions","wolfs","tigers");
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
