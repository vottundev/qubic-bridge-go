package assets

import "embed"

//go:embed *
var assets embed.FS

type SmartContract struct {
	Abi []byte
	Bin []byte
}

var Vottun1155 SmartContract = SmartContract{Abi: mustLoadBinaryFile("Vottun1155.abi"), Bin: mustLoadBinaryFile("Vottun1155.bin")}

// loads a file by name, if there is an error, panics
func mustLoadBinaryFile(fileName string) []byte {
	f, err := assets.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return f
}
