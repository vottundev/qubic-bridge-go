package assets

import "embed"

//go:embed *
var assets embed.FS

var TestEvent []byte = mustLoadBinaryFile("TestEvent.abi")

// loads a file by name, if there is an error, panics
func mustLoadBinaryFile(fileName string) []byte {
	f, err := assets.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return f
}
