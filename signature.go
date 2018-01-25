package phargo

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"io/ioutil"
	"os"
)

type signature struct {
	options Options
}

func (s *signature) check(filename string, f io.Reader) error {
	rest, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New("can't read rest of the file: " + err.Error())
	}

	return s.parse(filename, rest)
}

func (s *signature) parse(filename string, rest []byte) error {
	rLen := len(rest)
	if rLen < 8 {
		return errors.New("unexpected end of file, can't check last bytes")
	}

	if "GBMB" != string(rest[rLen-4:rLen]) {
		return errors.New("can't find GBMB constant at the end")
	}

	const (
		sigMD5    = 0x0001 //PHAR_SIG_MD5
		sigSHA1   = 0x0002 //PHAR_SIG_SHA1
		sigSHA256 = 0x0003 //PHAR_SIG_SHA256
	)

	f, _ := os.Open(filename)
	defer func() {
		_ = f.Close()
	}()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	//FILE_CONTENT + SIGNATURE + SIG_LENGTH + GBMB
	//              |<--      restBytes      --->|
	sigFlag := binary.LittleEndian.Uint32(rest[rLen-8 : rLen-4])
	var hasher hash.Hash
	var signatureLength int64
	algorithm := "UNKNOWN"

	switch sigFlag {
	case sigMD5:
		signatureLength = 16
		hasher = md5.New()
		algorithm = "MD5"

	case sigSHA1:
		signatureLength = 20
		hasher = sha1.New()
		algorithm = "SHA1"

	case sigSHA256:
		signatureLength = 32
		hasher = sha256.New()
		algorithm = "SHA256"

	default:
		return nil
	}

	if _, err := io.CopyN(hasher, f, stat.Size()-signatureLength-8); err != nil {
		return err
	}

	if !bytes.Equal(hasher.Sum(nil), rest[:signatureLength]) {
		return errors.New(algorithm + " hash of file is incorrect")
	}

	return nil
}
