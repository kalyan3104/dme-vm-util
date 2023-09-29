package mandosjsontest

import (
	"testing"

	mjparse "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/parse"
	mjwrite "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/write"
	"github.com/stretchr/testify/require"
)

func TestWriteScenario(t *testing.T) {
	contents, err := loadExampleFile("example.scen.json")
	require.Nil(t, err)

	p := mjparse.Parser{
		FileResolver: mjparse.NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}

	scenario, parseErr := p.ParseScenarioFile(contents)
	require.Nil(t, parseErr)

	serialized := mjwrite.ScenarioToJSONString(scenario)

	// good for debugging:
	// _ = ioutil.WriteFile("example_re.scen.json", []byte(serialized), 0644)

	require.Equal(t, contents, []byte(serialized))
}
