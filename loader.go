package gltf

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func FromBytes(data []byte) (*GlTF, error) {
	var root GlTF
	err := json.Unmarshal(data, &root)
	return &root, err
}

func FromFile(f *os.File) (*GlTF, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, errors.Join(errors.New("Could not Stat() file"), err)
	}

	if stat.Size() > 1<<30 {
		return nil, errors.New("File size is greater than 1GB soft limit")
	}

	b := make([]byte, stat.Size())
	if _, err = f.Read(b); err != nil {
		return nil, errors.Join(errors.New("Could not read file contents"), err)
	}

	gltf, err := FromBytes(b)
	if err != nil {
		return nil, errors.Join(errors.New("Failure parsing file as JSON"), err)
	}

	return gltf, nil
}

func FromFilename(name string) (*GlTF, error) {
	if f, err := os.Open(name); err != nil {
		return nil, err
	} else if root, err := FromFile(f); err != nil {
		return nil, err
	} else {
		root.meta.defaultSearchPath = filepath.Dir(name)
		return root, nil
	}
}
