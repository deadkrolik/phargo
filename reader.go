package phargo

import (
	"os"
)

//Reader - PHAR file parser
type Reader struct {
	options Options
}

//NewReader - creates parser with default options
func NewReader() *Reader {
	return &Reader{
		options: Options{
			MaxMetaDataLength: 10000,
			MaxManifestLength: 1048576 * 100,
			MaxFileNameLength: 1000,
			MaxAliasLength:    1000,
		},
	}
}

//SetOptions - applies options to parser
func (r *Reader) SetOptions(o Options) {
	r.options = o
}

//Parse - start parsing PHAR file
func (r *Reader) Parse(filename string) (File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return File{}, err
	}
	defer func() {
		_ = f.Close()
	}()

	var result File

	manifest := &manifest{options: r.options}
	offset, err := manifest.getOffset(f, 200, "__HALT_COMPILER(); ?>")
	if err != nil {
		return File{}, err
	}

	_, err = f.Seek(offset, 0)
	if err != nil {
		return File{}, err
	}

	err = manifest.parse(f)
	if err != nil {
		return File{}, err
	}
	result.Alias = string(manifest.Alias)
	result.Metadata = manifest.MetaSerialized
	result.Version = manifest.Version

	//files descriptions
	files := &files{options: r.options}
	result.Files, err = files.parse(f, manifest.EntitiesCount)
	if err != nil {
		return File{}, err
	}

	//check signature
	signature := &signature{options: r.options}
	err = signature.check(filename)
	if err != nil {
		return File{}, err
	}

	return result, nil
}
