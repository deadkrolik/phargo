package phargo

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type files struct {
	options Options
}

type entry struct {
	Name             string
	Timestamp        int64
	Size             int64
	Flags            uint32
	SizeUncompressed uint32
	SizeCompressed   uint32
	CRC              uint32
	MetaSerialized   []byte
}

func (f *files) parse(in io.Reader, count uint32) ([]PHARFile, error) {
	var i uint32
	var entries []entry
	var result []PHARFile

	for i = 0; i < count; i++ {
		entry, err := f.parseEntryHeader(in)
		if err != nil {
			return []PHARFile{}, err
		}
		entries = append(entries, entry)
	}

	//files data
	for _, entry := range entries {
		data, err := f.parseEntryData(in, &entry)
		if err != nil {
			return []PHARFile{}, err
		}

		result = append(result, PHARFile{
			Name:      entry.Name,
			Timestamp: entry.Timestamp,
			Metadata:  entry.MetaSerialized,
			Data:      data,
		})
	}

	return result, nil
}

func (f *files) parseEntryData(in io.Reader, entry *entry) ([]byte, error) {
	const (
		isCompressed = 0xF000 //PHAR_ENT_COMPRESSION_MASK
	)
	var buffer []byte

	if entry.Flags&isCompressed > 0 {
		return []byte{}, errors.New("can't parse compressed entry: " + entry.Name)
	}

	buffer = make([]byte, entry.SizeUncompressed)
	n, err := in.Read(buffer)
	if err != nil || n != int(entry.SizeUncompressed) {
		return []byte{}, errors.New("can't read entry data: " + entry.Name)
	}

	crc32q := crc32.MakeTable(0xedb88320)
	if entry.CRC != crc32.Checksum(buffer, crc32q) {
		return []byte{}, errors.New("entry has bad CRC: " + entry.Name)
	}

	return buffer, nil
}

func (f *files) parseEntryHeader(in io.Reader) (entry, error) {
	var e entry

	//length of entry name
	var nameLength uint32
	err := binary.Read(in, binary.LittleEndian, &nameLength)
	if err != nil || nameLength > f.options.MaxFileNameLength || nameLength == 0 {
		return entry{}, errors.New("can't read entry name length or it's too big")
	}

	//entry name
	buffer := make([]byte, nameLength)
	n, err := in.Read(buffer)
	if err != nil || n != int(nameLength) {
		return entry{}, errors.New("can't read entry name")
	}
	e.Name = string(buffer)

	//main entry info
	type entryBinary struct {
		SizeUncompressed uint32
		Timestamp        uint32
		SizeCompressed   uint32
		CRC              uint32
		Flags            uint32
		MetaLength       uint32
	}
	var eb entryBinary

	err = binary.Read(in, binary.LittleEndian, &eb)
	if err != nil {
		return entry{}, errors.New("can't read entry start: " + err.Error())
	}

	e.Timestamp = int64(eb.Timestamp)
	e.Size = int64(eb.SizeUncompressed)
	e.Flags = eb.Flags
	e.SizeCompressed = eb.SizeCompressed
	e.SizeUncompressed = eb.SizeUncompressed
	e.CRC = eb.CRC

	//read metadata
	if eb.MetaLength > f.options.MaxMetaDataLength {
		return entry{}, errors.New("entry metadata is too long: " + e.Name)
	}
	e.MetaSerialized = make([]byte, eb.MetaLength)

	n, err = in.Read(e.MetaSerialized)
	if err != nil || uint32(n) != eb.MetaLength {
		return entry{}, errors.New("can't read entry metadata: " + e.Name)
	}

	return e, nil
}
