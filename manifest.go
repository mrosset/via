package via

type File struct {
	Path string
	Type FileType
}

type FileType int

const (
	TypeFile FileType = iota
	TypeDir
	TypeLink
)

type Manifest struct {
	Plan  *Plan
	Files []*File
}
