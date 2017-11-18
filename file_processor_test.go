package stealer

import "testing"

func TestConvertFileName(t *testing.T) {
	testCases := []struct {
		filePath string
		expected string
	}{
		{
			filePath: "some/path/here.php",
			expected: "here.go",
		},

		{
			filePath: "some/path/here.java",
			expected: "",
		},
		{
			filePath: "some/error/path/here",
			expected: "",
		},
		{
			filePath: "some/error/path/here/",
			expected: "",
		},
	}

	for index, testCase := range testCases {
		result := convertFileName(testCase.filePath)
		if result != testCase.expected {
			t.Errorf("result = %+v , expected = %+v, at index = %+v\n", result, testCase.expected, index)
		}
	}
}
