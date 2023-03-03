# glTF

A simple file loader for glTF files in Go. This package is not intended to be used for manipulation or writing glTF
models to disk.

The entry point is any of FromBytes, FromFile, or FromFilename.  If using FromFilename (which is recomended), the
default search path for referenced files is automatically set to the same directory as the file.  

Once a file is loaded, calling Resolve will:

1) load binary data into byte slices stored in each Buffer
2) create sub-slice references to the bytes in each BufferView
3) replace all indexed references with a pointer to the resolved object; for example, a scene with a single node index
   in the source file will hold a pointer to the actual node after calling Resolve. That node will also have been
   resolved, so that you can walk the tree directly to the buffer data:
```go
    resolved := myGltf.Resolve(nil)
    var meshData []byte := resovled.Scene.Nodes[0].Mesh.Primitives.Attributes[gltf.POSITION].BufferView.Data
```

When using either of the first two forms, you will (likely) need to provide a search URI for referenced files, such as
textures or binary vertex data. Resolve accepts a slice of strings to allow searching multiple paths.

## Development Status

This package is working for loading of models and has partial support for cameras. It is not currently handling
textures, animations, etc. Features are being implemented in conjunction with development of
[https://github.com/bbredesen/gltf-viewer](a glTF model viewer) written in Go.

# License
MIT license. See the LICENSE file.

SPDX-License-Identifier: MIT

