package gltf

import (
	"testing"
)

type testData struct {
	Name               string
	ModelDirectoryPath string
}

var basicTests []testData = []testData{
	{"TriangleWithoutIndices", `.\glTF-Sample-Models\2.0\TriangleWithoutIndices\glTF\TriangleWithoutIndices.gltf`},
}

func TestBasic(t *testing.T) {
	// for _, tc := range basicTests {

	// }
}

func Test_TriangleWithoutIndices(t *testing.T) {
	if root, err := FromFilename(`.\glTF-Sample-Models\2.0\TriangleWithoutIndices\glTF\TriangleWithoutIndices.gltf`); err != nil {
		t.Fatalf("TriangleWithoutIndices: %v", err)
	} else {

		resolved, err := root.Resolve(nil)
		if err != nil {
			t.Fatalf(err.Error())
		}
		_ = resolved
		t.Log("Unmarshaled and resolved")
	}
}
