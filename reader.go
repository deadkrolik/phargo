package phargo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	haltCompiler = "__HALT_COMPILER(); ?>\r\n"
)

//Reader - PHAR file parser
type Reader struct {
	options Options
}

//NewReader - creates parser with default options
func NewReader() *Reader {
	r := Reader{}
	r.SetOptions(Options{
		MaxMetaLength:      10000,
		MaxManifestLength:  20000,
		MaxEntryNameLength: 1000,
	})

	return &r
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

	manifestOffset, err := r.getManifestOffset(f)
	if err != nil {
		return File{}, err
	}

	_, err = f.Seek(manifestOffset, 0)
	if err != nil {
		return File{}, err
	}

	manifest, err := r.parseManifest(f)
	if err != nil {
		return File{}, err
	}
	result.Metadata = manifest.MetaSerialized

	//files descriptions
	var i uint32
	var entries []entry
	for i = 0; i < manifest.EntitiesCount; i++ {
		entry, err := r.parseEntryHeader(f)
		if err != nil {
			return File{}, err
		}
		entries = append(entries, entry)
	}

	//files data
	for _, entry := range entries {
		data, err := r.parseEntryData(f, &entry)
		if err != nil {
			return File{}, err
		}

		result.Files = append(result.Files, PHARFile{
			Name:      entry.Name,
			Timestamp: entry.Timestamp,
			Data:      data,
		})
	}

	//check signature
	rest, err := ioutil.ReadAll(f)
	if err != nil {
		return File{}, errors.New("can't read rest of the file")
	}

	err = r.parseSignature(manifest, rest)
	if err != nil {
		return File{}, errors.New("can't parse signature")
	}

	return result, nil
}

func (r *Reader) parseSignature(m manifest, rest []byte) error {
	rLen := len(rest)
	if "GBMB" != string(rest[rLen-4:rLen]) {
		return errors.New("can't find GBMB constant at the end")
	}

	return nil
}

func (r *Reader) parseEntryData(f io.Reader, entry *entry) ([]byte, error) {
	const (
		isCompressed = 0xF000 //PHAR_ENT_COMPRESSION_MASK
	)
	var buffer []byte

	if entry.Flags&isCompressed > 0 {
		return []byte{}, errors.New("can't parse compressed entry: " + entry.Name)
	}

	buffer = make([]byte, entry.SizeUncompressed)
	n, err := f.Read(buffer)
	if err != nil || n != int(entry.SizeUncompressed) {
		return []byte{}, errors.New("can't read entry data: " + entry.Name)
	}

	crc32q := crc32.MakeTable(0xedb88320)
	if entry.CRC != crc32.Checksum(buffer, crc32q) {
		return []byte{}, errors.New("entry has bad CRC: " + entry.Name)
	}

	return buffer, nil
}

type entry struct {
	Name             string
	Timestamp        int64
	Size             int64
	Flags            uint32
	SizeUncompressed uint32
	SizeCompressed   uint32
	CRC              uint32
}

type entryBinary struct {
	SizeUncompressed uint32
	Timestamp        uint32
	SizeCompressed   uint32
	CRC              uint32
	Flags            uint32
	MetaLength       uint32
}

func (r *Reader) parseEntryHeader(f io.Reader) (entry, error) {
	var e entry

	//length of entry name
	var nameLength uint32
	err := binary.Read(f, binary.LittleEndian, &nameLength)
	if err != nil || nameLength > r.options.MaxEntryNameLength || nameLength == 0 {
		return entry{}, errors.New("can't read entry name length or it's too big")
	}

	//entry name
	buffer := make([]byte, nameLength)
	n, err := f.Read(buffer)
	if err != nil || n != int(nameLength) {
		return entry{}, errors.New("can't read entry name")
	}
	e.Name = string(buffer)

	//main entry info
	var eb entryBinary
	err = binary.Read(f, binary.LittleEndian, &eb)
	if err != nil {
		return entry{}, errors.New("can't read entry start: " + err.Error())
	}

	e.Timestamp = int64(eb.Timestamp)
	e.Size = int64(eb.SizeUncompressed)
	e.Flags = eb.Flags
	e.SizeCompressed = eb.SizeCompressed
	e.SizeUncompressed = eb.SizeUncompressed
	e.CRC = eb.CRC

	//metadata of entry
	_, err = io.CopyN(ioutil.Discard, f, int64(eb.MetaLength))
	if err != nil {
		return entry{}, errors.New("can't skip metadata of entry")
	}

	return e, nil
}

type manifest struct {
	Length         uint32
	EntitiesCount  uint32
	Version        string
	Flags          uint32
	AliasLength    uint32
	MetaLength     uint32
	MetaSerialized []byte
	IsSigned       bool
}

type manifestBinary struct {
	Length        uint32
	EntitiesCount uint32
	Version       uint16
	Flags         uint32
	AliasLength   uint32
}

func (r *Reader) parseManifest(f io.Reader) (manifest, error) {
	var (
		m  manifest
		mb manifestBinary
	)

	err := binary.Read(f, binary.LittleEndian, &mb)
	if err != nil {
		return manifest{}, errors.New("can't read manifest header: " + err.Error())
	}

	m.Length = mb.Length
	m.EntitiesCount = mb.EntitiesCount
	m.Version = fmt.Sprintf("%d.%d.%d", (mb.Version<<12)>>12, ((mb.Version>>4)<<12)>>12, ((mb.Version>>8)<<12)>>12)
	m.Flags = mb.Flags
	m.AliasLength = mb.AliasLength

	_, err = io.CopyN(ioutil.Discard, f, int64(m.AliasLength))
	if err != nil {
		return manifest{}, errors.New("can't skip alias while reading manifest")
	}

	err = binary.Read(f, binary.LittleEndian, &m.MetaLength)
	if err != nil {
		return manifest{}, errors.New("can't read manifest metadata length")
	}

	if m.MetaLength > r.options.MaxMetaLength {
		return manifest{}, errors.New("metadata length of manifest more than MaxMetaLength")
	}
	m.MetaSerialized = make([]byte, m.MetaLength)

	n, err := f.Read(m.MetaSerialized)
	if err != nil || uint32(n) != m.MetaLength {
		return manifest{}, errors.New("can't read manifest metadata")
	}

	m.IsSigned = m.Flags&0x10000 > 0 //PHAR_HDR_SIGNATURE

	return m, nil
}

func (r *Reader) getManifestOffset(f io.Reader) (int64, error) {
	buffer := make([]byte, 200)
	before := make([]byte, 200)
	var position int64

	for {
		n, err := f.Read(buffer)
		if err == io.EOF {
			return 0, errors.New("unexpected EOF while looking for manifest")
		}

		search := append(before, buffer...)
		index := strings.Index(string(search), haltCompiler)

		if index >= 0 {
			offset := position + int64(index) - int64(n) + int64(len(haltCompiler))
			if offset > r.options.MaxManifestLength {
				return 0, errors.New("manifest length more than MaxManifestLength")
			}

			return offset, nil
		}

		position += int64(n)
	}
}

//SetOptions - applies options to parser
func (r *Reader) SetOptions(o Options) {
	r.options = o
}
