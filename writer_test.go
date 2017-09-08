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
			expected: "someVar := []int{1,2,3}",
		},

		{
			varName:  "someVar",
			values:   []string{"true", "false", "true", "false", "true"},
			expected: "someVar := []bool{true,false,true,false,true}",
		},
		{
			varName: "someVar",
			values:  []string{"hello", "world", "bray", "wkwk", "hehe"},

			expected: "someVar := []string{\"hello\",\"world\",\"bray\",\"wkwk\",\"hehe\"}",
		},
	}

	for _, testObject := range testObjects {
		actual := variableFormat(testObject.varName, testObject.values)
		if actual != testObject.expected {
			t.Errorf("actual = %+v, expected = %+v\n", actual, testObject.expected)
		}
	}
}
