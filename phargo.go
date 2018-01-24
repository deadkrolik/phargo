package phargo

//Options - parser options
type Options struct {
	MaxMetaDataLength uint32
	MaxManifestLength int64
	MaxFileNameLength uint32
}

//File - parsed PHAR-file
type File struct {
	Version  string
	Alias    string
	Metadata []byte
	Files    []PHARFile
}

//PHARFile - file inside PHAR-archive
type PHARFile struct {
	Name      string
	Timestamp int64
	Metadata  []byte
	Data      []byte
}
