package server

import (
	"testing"
)

func TestParse(t *testing.T) {
	tt := []struct{
		name string
		in string
	}{
		{"CONCAT with strings only", `CONCAT("Hello ", "World")`},
		{"CONCAT with strings and float", `CONCAT("Hello ", 3.1416)`},
		{"CONCAT with strings and int", `CONCAT("Hello ", 42)`},
		{"EXACT function", `EXACT("Hello", "Hello")`},
		{"FIND function with start position", `FIND("l", "Hello", 2)`},
		{"FIND function without start position", `FIND("l", "Hello")`},
		{"LEFT with string and no num chars", `LEFT("Hello World")`},
		{"LEFT with string and with num chars", `LEFT("Hello World", 5)`},
	}
	var errCnt int
	for _, tu := range tt {
		_, err := parse(tu.in)
		if err != nil {
			t.Logf("Test Case: %v, Error on Parse: %v", tu.name, err.Error())
			errCnt++
		}
	}
	if errCnt > 0 {
		t.Errorf("Failed %v of %v cases", errCnt, len(tt))
	}
}

func TestBuildCache(t *testing.T) {
	t.Errorf("Need to implement test")
}
