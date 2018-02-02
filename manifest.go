package phargo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

type manifest struct {
	Length         uint32
	EntitiesCount  uint32
	Version        string
	Flags          uint32
	Alias          []byte
	AliasLength    uint32
	MetaLength     uint32
	MetaSerialized []byte
	IsSigned       bool

	options Options
}

func (m *manifest) parse(f io.Reader) error {
	type mBinary struct {
		Length        uint32
		EntitiesCount uint32
		Version       uint16
		Flags         uint32
		AliasLength   uint32
	}
	var mb mBinary

	err := binary.Read(f, binary.LittleEndian, &mb)
	if err != nil {
		return errors.New("can't read manifest header: " + err.Error())
	}

	m.Length = mb.Length
	m.EntitiesCount = mb.EntitiesCount
	m.Version = fmt.Sprintf("%d.%d.%d", (mb.Version<<12)>>12, ((mb.Version>>4)<<12)>>12, ((mb.Version>>8)<<12)>>12)
	m.Flags = mb.Flags
	m.AliasLength = mb.AliasLength

	if m.AliasLength > m.options.MaxAliasLength {
		return errors.New("manifest alias length more than MaxAlasLength")
	}

	m.Alias = make([]byte, m.AliasLength)
	n, err := f.Read(m.Alias)
	if err != nil || uint32(n) != m.AliasLength {
		return errors.New("can't read manifest alias")
	}

	err = binary.Read(f, binary.LittleEndian, &m.MetaLength)
	if err != nil {
		return errors.New("can't read manifest metadata length: " + err.Error())
	}

	if m.MetaLength > m.options.MaxMetaDataLength {
		return errors.New("metadata length of manifest more than MaxMetaDataLength")
	}
	m.MetaSerialized = make([]byte, m.MetaLength)

	n, err = f.Read(m.MetaSerialized)
	if err != nil || uint32(n) != m.MetaLength {
		return errors.New("can't read manifest metadata")
	}

	m.IsSigned = m.Flags&0x10000 > 0 //PHAR_HDR_SIGNATURE
	return nil
}

func (m *manifest) getOffset(f io.Reader, bufSize int64, haltCompiler string) (int64, error) {
	buffer := make([]byte, bufSize)
	before := make([]byte, bufSize)
	var position int64

	for {
		n, err := f.Read(buffer)
		if err != nil {
			return 0, errors.New("can't find haltCompiler: " + err.Error())
		}

		search := append(before, buffer...)
		index := strings.Index(string(search), haltCompiler)

		if index >= 0 {
			offset := position + int64(index) - bufSize + int64(len(haltCompiler))
			if index+len(haltCompiler) >= len(search) {
				return 0, errors.New("unexpected end of file")
			}

			//optional \r\n or \n
			var nextChar = search[index+len(haltCompiler)]
			var nextNextChar = search[index+len(haltCompiler)+1]
			if nextChar == '\r' && nextNextChar == '\n' {
				offset += 2
			}
			if nextChar == '\n' {
				offset++
			}

			return offset, nil
		}

		position += int64(n)
		copy(before, buffer)

		if position > m.options.MaxManifestLength {
			return 0, errors.New("manifest length more than MaxManifestLength")
		}
	}
}
