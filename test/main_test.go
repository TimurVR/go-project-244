package code

import (
	"code"
	"reflect"
	"testing"
)

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name     string
		file1    string
		file2    string
		expected string
	}{
		{
			name:     "Empty files",
			file1:    "../testdata/file1_empty.json",
			file2:    "../testdata/file2_empty.json",
			expected: "\n{\n}\n",
		},
		{
			name:     "Simple files with some differences",
			file1:    "../testdata/file1_simple.json",
			file2:    "../testdata/file2_simple.json",
			expected: "\n{\n    a: 1\n  - b: test\n  + b: test2\n  - c: true\n  + d: false\n}\n",
		},
		{
			name:     "Identical simple files",
			file1:    "../testdata/file1_simple.json",
			file2:    "../testdata/file1_simple.json",
			expected: "\n{\n    a: 1\n    b: test\n    c: true\n}\n",
		},
		{
			name:     "First file empty, second with data",
			file1:    "../testdata/file1_empty.json",
			file2:    "../testdata/file1_simple.json",
			expected: "\n{\n  + a: 1\n  + b: test\n  + c: true\n}\n",
		},
		{
			name:     "First file with data, second empty",
			file1:    "../testdata/file1_simple.json",
			file2:    "../testdata/file1_empty.json",
			expected: "\n{\n  - a: 1\n  - b: test\n  - c: true\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			map1 := code.Parsing(tt.file1)
			map2 := code.Parsing(tt.file2)
			result := code.GenDiff(map1, map2)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GenDiff() = %v, want %v", result, tt.expected)
			}
		})
	}
}