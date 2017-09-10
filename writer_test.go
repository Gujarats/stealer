package stealer

import "testing"

func TestVaraibleFormat(t *testing.T) {
	testObjects := []struct {
		varName  string
		values   []string
		expected string
	}{
		{
			varName:  "someVar",
			values:   []string{"1", "2", "3"},
			expected: "var someVar = []int{1,2,3}",
		},

		{
			varName:  "someVar",
			values:   []string{"true", "false", "true", "false", "true"},
			expected: "var someVar = []bool{true,false,true,false,true}",
		},
		{
			varName: "someVar",
			values:  []string{"hello", "world", "bray", "wkwk", "hehe"},

			expected: "var someVar = []string{\"hello\",\"world\",\"bray\",\"wkwk\",\"hehe\"}",
		},

		// test item for more than 10
		{
			varName: "someVar",
			values:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},

			expected: "var someVar = []int{1,2,3,4,5,6,7,8,9,10,\n11,12}",
		},
	}

	for _, testObject := range testObjects {
		actual := variableFormat(testObject.varName, testObject.values)
		if actual != testObject.expected {
			t.Errorf("actual = %+v, expected = %+v\n", actual, testObject.expected)
		}
	}
}
