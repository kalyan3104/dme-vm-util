package mandosjsontest

import (
	"io/ioutil"
	"os"
	"testing"

	mjparse "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/parse"
	mjwrite "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/write"
	"github.com/stretchr/testify/require"
)

func loadExampleFile(path string) ([]byte, error) {
	// Open our jsonFile
	var jsonFile *os.File
	var err error
	jsonFile, err = os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	return ioutil.ReadAll(jsonFile)
}

func TestWriteTest(t *testing.T) {
	contents, err := loadExampleFile("example.test.json")
	require.Nil(t, err)

	p := mjparse.Parser{
		FileResolver: mjparse.NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}

	testTopLevel, parseErr := p.ParseTestFile(contents)
	require.Nil(t, parseErr)

	serialized := mjwrite.TestToJSONString(testTopLevel)

	// good for debugging:
	// _ = ioutil.WriteFile("example.re__.json", []byte(serialized), 0644)

	require.Equal(t, contents, []byte(serialized))
}
