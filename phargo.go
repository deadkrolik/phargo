package phargo

//Options - parser options
type Options struct {
	MaxMetaLength      uint32
	MaxManifestLength  int64
	MaxEntryNameLength uint32
}

//File - parsed PHAR-file
type File struct {
	Metadata []byte
	Files    []PHARFile
}

//PHARFile - file inside PHAR-archive
type PHARFile struct {
	Name      string
	Timestamp int64
	Data      []byte
}
