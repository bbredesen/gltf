package gltf

// ResolvedGlTF is a GlTF with all indexed node references and binary data resolved and loaded. The source GlTF struct
// is included as a pointer, so all of the source data can be referenced as neccessary. Fields for resolved data types
// are named the same as the sources in a GlTF struct, masking the original data. The original data can still be
// accessed through the embedded GlTF pointer:
// ```
// var resolved *ResolvedGlTF
// ...
// resolved.Nodes => slice of ResolvedNode
// resolved.GlTF.Nodes => slice of Node
// ```
type ResolvedGlTF struct {
	*GlTF

	Animations  []ResolvedAnimation
	Buffers     []ResolvedBuffer
	BufferViews []ResolvedBufferView
	Accessors   []ResolvedAccessor
	Materials   []ResolvedMaterial
	Meshes      []ResolvedMesh
	Nodes       []ResolvedNode

	Scene  *ResolvedScene
	Scenes []ResolvedScene
}

type ResolvedNode struct {
	*Node
	Children []*ResolvedNode
	Mesh     *ResolvedMesh
}

type ResolvedScene struct {
	*Scene
	Nodes []*ResolvedNode
}

type ResolvedBuffer struct {
	*Buffer
	Data []byte
}

type ResolvedBufferView struct {
	*BufferView
	Buffer *ResolvedBuffer
	// Data in the BufferView is a subslice (i.e. shared memory) from the source buffer data
	Data []byte
}

type ResolvedAccessor struct {
	*Accessor
	BufferView *ResolvedBufferView
}

type ResolvedMesh struct {
	*Mesh
	Primitives []ResolvedPrimitive
}

type ResolvedPrimitive struct {
	*Primitive
	Attributes map[AttributeKey]*ResolvedAccessor
	Indices    *ResolvedAccessor
	Material   *ResolvedMaterial
}

type ResolvedMaterial struct {
	*Material

	// Camera *ResolvedCamera
	// TODO...
}

type ResolvedAnimation struct {
	*Animation
	Channels []ResolvedAnimationChannel
	Samplers []ResolvedAnimationSampler
}

type ResolvedAnimationChannel struct {
	*AnimationChannel
	Sampler *ResolvedAnimationSampler
	Target  ResolvedAnimationChannelTarget
}

type ResolvedAnimationChannelTarget struct {
	*AnimationChannelTarget
	Node *ResolvedNode
}

type ResolvedAnimationSampler struct {
	*AnimationSampler
	Input  *ResolvedAccessor
	Output *ResolvedAccessor
}
