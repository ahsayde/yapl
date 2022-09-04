package test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ahsayde/yapl/yapl"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type TestSuite struct {
	ID        string                 `yaml:"id"`
	Policy    map[string]interface{} `yaml:"policy"`
	Testcases []TestCase             `yaml:"tests"`
}

type TestCase struct {
	ID     string                 `yaml:"id"`
	Input  map[string]interface{} `yaml:"input"`
	Params map[string]interface{} `yaml:"params"`
	Result *yapl.Result           `yaml:"result"`
	Errors []string               `yaml:"errors"`
}

func readTestFile(path string) ([]TestSuite, error) {
	var testsuites []TestSuite

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(raw))

	for {
		var testsuite TestSuite
		if err := decoder.Decode(&testsuite); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		testsuites = append(testsuites, testsuite)
	}

	return testsuites, nil
}

func Test(t *testing.T) {
	files, err := filepath.Glob("./testcases/*yaml")
	if err != nil {
		t.Fatalf(err.Error())
	}
	var testsuites []TestSuite
	for i := range files {
		tests, err := readTestFile(files[i])
		if err != nil {
			t.Fatal(err.Error())
		}
		testsuites = append(testsuites, tests...)
	}
	for _, testsuite := range testsuites {
		for _, testcase := range testsuite.Testcases {
			testID := fmt.Sprintf("%s.%s", testsuite.ID, testcase.ID)
			raw, err := yaml.Marshal(testsuite.Policy)
			if err != nil {
				t.Fatal(err)
			}
			policy, err := yapl.Parse(raw)
			if err != nil {
				t.Fatal(err)
			}
			result, err := policy.Eval(testcase.Input, testcase.Params)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testcase.Result, result, "testcase: %s failed", testID)
		}
	}
}
