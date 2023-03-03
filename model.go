package gltf

import (
	"fmt"

	"github.com/bbredesen/vkm"
)

// Spec: The root object for a glTF asset.
type GlTF struct {
	// Asset is a required field
	Asset Asset `json:"asset"`

	ExtensionsUsed     []string `json:"extensionsUsed"`
	ExtensionsRequired []string `json:"extensionsRequired"`

	Accessors   []Accessor   `json:"accessors"`
	Animations  []Animation  `json:"animations"`
	Buffers     []Buffer     `json:"buffers"`
	BufferViews []BufferView `json:"bufferViews"`
	Cameras     []Camera     `json:"cameras"`
	Images      []Image      `json:"images"`
	Materials   []Material   `json:"materials"`
	Meshes      []Mesh       `json:"meshes"`
	Nodes       []Node       `json:"nodes"`
	Samplers    []Sampler    `json:"samplers"`
	Scene       *uint        `json:"scene,omitempty"` // Spec: Scene is an optional reference to the default scene for this asset, as an index in the Scenes array.
	Scenes      []Scene      `json:"scenes"`
	Skins       []Skin       `json:"skins"`
	Textures    []Texture    `json:"textures"`

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`

	meta struct {
		defaultSearchPath string
	}
}

// Spec: Metadata about the glTF asset.
type Asset struct {
	// Spec: A copyright message suitable for display to credit the content creator.
	Copyright string `json:"copyright,omitempty"`
	// Spec: Tool that generated this glTF model.  Useful for debugging.
	Generator string `json:"generator,omitempty"`
	// Spec: The glTF version in the form of `<major>.<minor>` that this asset targets.
	Version string `json:"version"`
	// The minimum glTF version in the form of `<major>.<minor>` that this asset targets. This property **MUST NOT** be greater than the asset version.
	MinVersion string `json:"minVersion,omitempty"`
}

type Scene struct {
	Nodes []uint `json:"nodes"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type GlTFId string

type Node struct {
	Camera      *int        `json:"camera"`
	Children    []uint      `json:"children"`
	Skin        *int        `json:"skin"`
	Matrix      [16]float32 `json:"matrix"`
	Mesh        *uint       `json:"mesh,omitempty"`
	Rotation    *vkm.Vec    `json:"rotation"`
	Scale       *vkm.Vec3   `json:"scale"`
	Translation *vkm.Vec3   `json:"translation"`
	Weights     []float32   `json:"weights"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type Material struct {
	// TODO

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type Mesh struct {
	Primitives []Primitive `json:"primitives"`
	Weights    []float32   `json:"weights"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type Primitive struct {
	Attributes map[AttributeKey]int `json:"attributes"`
	Indices    *uint                `json:"indices,omitempty"`
	// May be null, indicating "default" material
	Material *int      `json:"material,omitempty"`
	Mode     *ModeEnum `json:"mode,omitempty"`
	Targets  []string  `json:"targets"` // TODO

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

// Constant strings representing standard attributes, i.e. what implementations should support per the spec. Note that
// indexes for TEXCOORD_n, COLOR_n, JOINTS_n, and WEIGHTS_n can go to 9.

type AttributeKey string

const (
	POSITION   AttributeKey = "POSITION"
	NORMAL                  = "NORMAL"
	TANGENT                 = "TANGENT"
	TEXCOORD_0              = "TEXCOORD_0"
	TEXCOORD_1              = "TEXCOORD_1"
	COLOR_0                 = "COLOR_0"
	JOINTS_0                = "JOINTS_0"
	WEIGHTS_0               = "WEIGHTS_0"
)

type ModeEnum int

const (
	POINTS ModeEnum = iota
	LINES
	LINE_LOOP
	LINE_STRIP
	TRIANGLES
	TRIANGLE_STRIP
	TRIANGLE_FAN
)

// These are supposed to be generic maps
type Extensions map[string]any
type Extras map[string]any

// Accessor: see https://registry.khronos.org/glTF/specs/2.0/glTF-2.0.html#schema-reference-accessor
type Accessor struct {
	BufferView    uint              `json:"bufferView"`
	ByteOffset    uint              `json:"byteOffset"`
	ComponentType ComponentTypeEnum `json:"componentType"`
	Normalized    bool              `json:"normalized"`
	Count         int               `json:"count"`
	Type          AccessorTypeEnum  `json:"type"`
	Max           []float64         `json:"max"`
	Min           []float64         `json:"min"`
	Sparse        SparseAccessor    `json:"sparse"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type ComponentTypeEnum int

const (
	BYTE           ComponentTypeEnum = 5120
	UNSIGNED_BYTE  ComponentTypeEnum = 5121
	SHORT          ComponentTypeEnum = 5122
	UNSIGNED_SHORT ComponentTypeEnum = 5123
	UNSIGNED_INT   ComponentTypeEnum = 5125
	FLOAT          ComponentTypeEnum = 5126
)

// Size returns the byte size of the component specified by an Accessor. Note that the size of component types are
// determined by the glTF spec, not the machine that the file is being interpreted by. This function will panic() if
// called on an enum value not defined by this module.
func (c ComponentTypeEnum) Size() int {
	switch c {
	case BYTE:
		fallthrough
	case UNSIGNED_BYTE:
		return 1
	case SHORT:
		fallthrough
	case UNSIGNED_SHORT:
		return 2
	case UNSIGNED_INT:
		fallthrough
	case FLOAT:
		return 4
	}
	panic(fmt.Sprint("unknown ComponentType:", c))
}

type AccessorTypeEnum string

const (
	SCALAR AccessorTypeEnum = "SCALAR"
	VEC2   AccessorTypeEnum = "VEC2"
	VEC3   AccessorTypeEnum = "VEC3"
	VEC4   AccessorTypeEnum = "VEC4"
	MAT2   AccessorTypeEnum = "MAT2"
	MAT3   AccessorTypeEnum = "MAT3"
	MAT4   AccessorTypeEnum = "MAT4"
)

// Count returns the number of individual components, without regard to the byte size of that component, in each element
// as defined by an Accessor. This function will panic() if called on an enum value not defined by this module.
func (ate AccessorTypeEnum) Count() int {
	switch ate {
	case SCALAR:
		return 1
	case VEC2:
		return 2
	case VEC3:
		return 3
	case VEC4:
		return 4
	case MAT2:
		return 4
	case MAT3:
		return 9
	case MAT4:
		return 16
	}
	panic(fmt.Sprint("unknown AccessorType:", ate))
}

// Stride is a convenience function returning the number of bytes in each element as defined by this Accessor.
func (a *Accessor) Stride() int {
	return a.ComponentType.Size() * a.Type.Count()
}

type Buffer struct {
	Uri        string `json:"uri,omitempty"`
	ByteLength uint   `json:"byteLength"`

	Name       GlTFId `json:"name"`
	Extensions `json:"extensions"`
	Extras     `json:"extras"`
}

type BufferView struct {
	Buffer     uint             `json:"buffer"`
	ByteOffset uint             `json:"byteOffset"`
	ByteLength uint             `json:"byteLength"`
	ByteStride uint             `json:"byteStride"`
	Target     BufferTargetEnum `json:"target"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type BufferTargetEnum int

const (
	ARRAY_BUFFER         BufferTargetEnum = 34962
	ELEMENT_ARRAY_BUFFER BufferTargetEnum = 34963
)

type Camera struct {
	Orthographic CameraOrthographic `json:"orthographic,omitempty"`
	Perspective  CameraPerpsective  `json:"perspective,omitempty"`
	Type         CameraType         `json:"type"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type CameraOrthographic struct {
	Xmag  float32 `json:"xmag"`
	Ymag  float32 `json:"ymag"`
	Zfar  float32 `json:"zfar"`
	Znear float32 `json:"znear"`

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type CameraPerpsective struct {
	AspectRatio float32 `json:"aspectRatio"`
	Yfov        float32 `json:"yfov"`
	Zfar        float32 `json:"zfar"`
	Znear       float32 `json:"znear"`

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type CameraType string

const (
	PERSPECTIVE  CameraType = "perspective"
	ORTHOGRAPHIC CameraType = "orthographic"
)

type Animation struct {
	Channels []AnimationChannel `json:"channels"`
	Samplers []AnimationSampler `json:"samplers"`

	Name       GlTFId `json:"name,omitempty"`
	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type AnimationChannel struct {
	Sampler uint                   `json:"sampler"`
	Target  AnimationChannelTarget `json:"target"`

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type AnimationChannelTarget struct {
	Node *uint                      `json:"node,omitempty"`
	Path AnimationChannelTargetPath `json:"path"`

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type AnimationChannelTargetPath string

const (
	TRANSLATION AnimationChannelTargetPath = "translation"
	ROTATION    AnimationChannelTargetPath = "rotation"
	SCALE       AnimationChannelTargetPath = "scale"
	WEIGHTS     AnimationChannelTargetPath = "weights"
)

type AnimationSampler struct {
	Input         int                           `json:"input"`
	Interpolation AnimationSamplerInterpolation `json:"interpolation,omitempty"`
	Output        int                           `json:"output"`

	Extensions `json:"extensions,omitempty"`
	Extras     `json:"extras,omitempty"`
}

type AnimationSamplerInterpolation string

const (
	LINEAR       AnimationSamplerInterpolation = "linear"
	STEP         AnimationSamplerInterpolation = "step"
	CUBIC_SPLINE AnimationSamplerInterpolation = "cubicspline"
)

// Everything below here is TODO

type Image map[string]any
type Sampler map[string]any
type Skin map[string]any
type Texture map[string]any
type SparseAccessor map[string]any
