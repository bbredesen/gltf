package gltf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// TODO: resolve's probably need error checking for indexing issues...e.g. scene references a node index that doesn't
// exist. This is a problem with the source being malformed, but is easy enough to check here.

// Resolve processes a glTF structure and returns a ResolvedGlTF that can be used in a program without further
// processing. The returned value will have:
//
// * All indexed resource references resolved as pointers to that resource. For example, an array of Node indices will
// be converted to pointers to the actual Nodes.
// * All binary resources will be loaded into their Buffer objects, using the provided uriSearchPaths for name
// resolution. Note that the binary load will
// not interpret or validate the data, it will simply be loaded as-is into the buffer. An error will be raised if the binary data provided is
// shorter than the Buffer's ByteLength. NOTE HOWEVER, that the data will still be loaded into the Buffer and zero
// padded. The any data at the URI beyond the Buffer's ByteLength will be ignored.
//
// Paths will be searched in the order provided and the first matching file will be used. If the GlTF instance  was loaded from a
// file object, or from a file name, and if uriSearchPath is empty, then that location will be searched by default.
func (gltf *GlTF) Resolve(uriSearchPath []string) (*ResolvedGlTF, error) {
	rval := &ResolvedGlTF{GlTF: gltf}

	for i := range gltf.Buffers {
		if rb, err := gltf.Buffers[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Buffers = append(rval.Buffers, rb)
		}
	}

	for i := range gltf.BufferViews {
		// tmp := bv
		if rbv, err := gltf.BufferViews[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.BufferViews = append(rval.BufferViews, rbv)
		}
	}

	for i := range gltf.Accessors {
		if rac, err := gltf.Accessors[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Accessors = append(rval.Accessors, rac)
		}
	}

	for i := range gltf.Materials {
		if rm, err := gltf.Materials[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Materials = append(rval.Materials, rm)
		}
	}

	for i := range gltf.Meshes {
		if rm, err := gltf.Meshes[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Meshes = append(rval.Meshes, rm)
		}
	}

	for i := range gltf.Nodes {
		if rn, err := gltf.Nodes[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Nodes = append(rval.Nodes, rn)
		}
	}
	for i := range rval.Nodes {
		rval.Nodes[i].populate(rval)
	}

	for i := range gltf.Animations {
		if ra, err := gltf.Animations[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Animations = append(rval.Animations, ra)
		}
	}

	for i := range gltf.Scenes {
		if rs, err := gltf.Scenes[i].resolve(rval); err != nil {
			return rval, err
		} else {
			rval.Scenes = append(rval.Scenes, rs)
		}
	}
	if gltf.Scene != nil {
		rval.Scene = &rval.Scenes[*gltf.Scene]
	}

	// mesh, accessor, buffer view, buffer
	// for _, m := range gltf.Meshes {
	// 	if rm, err := m.populate(rval); err != nil {
	// 		return rval, err
	// 	} else {
	// 		rval.Meshes = append(rval.Meshes, rm)
	// 	}
	// }
	return rval, nil
}

func (node *Node) resolve(root *ResolvedGlTF) (ResolvedNode, error) {
	rval := ResolvedNode{
		Node: node,
	}

	rval.Children = make([]*ResolvedNode, len(node.Children))

	if node.Mesh != nil {
		rval.Mesh = &root.Meshes[*node.Mesh]
	}

	return rval, nil
}

func (node *ResolvedNode) populate(root *ResolvedGlTF) {
	for i, childIdx := range node.Node.Children {
		node.Children[i] = &root.Nodes[childIdx]
	}
}

func (buf *Buffer) resolve(root *ResolvedGlTF) (ResolvedBuffer, error) {
	rval := ResolvedBuffer{
		Buffer: buf,
	}

	// Load the data; assume the URI just a file (for the moment)
	data, err := os.ReadFile(root.meta.defaultSearchPath + string(filepath.Separator) + buf.Uri)
	if err != nil {
		return rval, err
	}

	if len(data) < int(buf.ByteLength) {
		return rval, errors.New(fmt.Sprintf("Binary size was smaller than specification: %s, expected >= %d bytes, got %d bytes", buf.Uri, buf.ByteLength, len(data)))
	}

	rval.Data = data
	return rval, nil
}

func (bv *BufferView) resolve(root *ResolvedGlTF) (ResolvedBufferView, error) {
	rval := ResolvedBufferView{
		BufferView: bv,
	}

	rval.Buffer = &root.Buffers[bv.Buffer]
	rval.Data = rval.Buffer.Data[bv.ByteOffset : bv.ByteOffset+bv.ByteLength]

	return rval, nil
}

func (a *Accessor) resolve(root *ResolvedGlTF) (ResolvedAccessor, error) {
	rval := ResolvedAccessor{
		Accessor: a,
	}

	rval.BufferView = &root.BufferViews[a.BufferView]

	return rval, nil
}

func (s *Scene) resolve(root *ResolvedGlTF) (ResolvedScene, error) {
	rval := ResolvedScene{
		Scene: s,
	}

	rval.Nodes = make([]*ResolvedNode, len(s.Nodes))
	for i := range s.Nodes {
		rval.Nodes[i] = &root.Nodes[i]
	}

	return rval, nil
}

func (m *Material) resolve(root *ResolvedGlTF) (ResolvedMaterial, error) {
	rval := ResolvedMaterial{
		Material: m,
	}

	return rval, nil
}

func (m *Mesh) resolve(root *ResolvedGlTF) (ResolvedMesh, error) {
	rval := ResolvedMesh{
		Mesh: m,
	}

	for _, primitive := range m.Primitives {
		if rp, err := primitive.resolve(root); err != nil {
			return rval, err
		} else {
			rval.Primitives = append(rval.Primitives, rp)
		}
	}

	return rval, nil
}

func (p *Primitive) resolve(root *ResolvedGlTF) (ResolvedPrimitive, error) {
	rval := ResolvedPrimitive{
		Primitive: p,
	}

	if p.Material != nil {
		rval.Material = &root.Materials[*p.Material]
	}

	rval.Attributes = make(map[AttributeKey]*ResolvedAccessor, len(p.Attributes))
	for k, attrIdx := range p.Attributes {
		rval.Attributes[k] = &root.Accessors[attrIdx]
	}

	if p.Indices != nil {
		rval.Indices = &root.Accessors[*p.Indices]
	}

	//temp comment
	return rval, nil
}

func (a *Animation) resolve(root *ResolvedGlTF) (ResolvedAnimation, error) {
	rval := ResolvedAnimation{
		Animation: a,
	}

	rval.Samplers = make([]ResolvedAnimationSampler, len(a.Samplers))
	for i := range a.Samplers {
		if s, err := a.Samplers[i].resolve(root); err != nil {
			return rval, err
		} else {
			rval.Samplers[i] = s
		}
	}

	rval.Channels = make([]ResolvedAnimationChannel, len(a.Channels))
	for i := range a.Channels {
		if ch, err := a.Channels[i].resolve(root, rval.Samplers); err != nil {
			return rval, err
		} else {
			rval.Channels[i] = ch
		}
	}

	return rval, nil
}

func (as *AnimationSampler) resolve(root *ResolvedGlTF) (ResolvedAnimationSampler, error) {
	rval := ResolvedAnimationSampler{
		AnimationSampler: as,
	}

	rval.Input = &root.Accessors[as.Input]
	rval.Output = &root.Accessors[as.Output]

	return rval, nil
}

func (ac *AnimationChannel) resolve(root *ResolvedGlTF, samplers []ResolvedAnimationSampler) (ResolvedAnimationChannel, error) {
	rval := ResolvedAnimationChannel{
		AnimationChannel: ac,
	}

	var err error
	rval.Sampler = &samplers[ac.Sampler]
	rval.Target, err = ac.Target.resolve(root)
	return rval, err
}

func (act *AnimationChannelTarget) resolve(root *ResolvedGlTF) (ResolvedAnimationChannelTarget, error) {
	rval := ResolvedAnimationChannelTarget{
		AnimationChannelTarget: act,
	}

	if act.Node != nil {
		rval.Node = &root.Nodes[*act.Node]
	}

	return rval, nil
}
