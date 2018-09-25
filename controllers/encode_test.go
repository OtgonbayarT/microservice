package controllers

import "testing"
func TestEncode (t *testing.T) {		

	encodeTests := []struct {		
		input       string
        expected    uint32
	}{
		{	
			input: 			"https://www.youtube.com/watch?v=PDxcEzu62jk",
			expected:		2359611651,
		},
		{	
			input: 			"",
			expected:		2166136261,
		},

	}
	for _, test := range encodeTests {
		actual := Hash(test.input)
		if actual != test.expected {
				t.Errorf("String(%s): expected %d, actual %d", test.input, test.expected, actual)
		}
}
}