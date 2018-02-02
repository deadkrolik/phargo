package phargo

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"os"
)

type signature struct {
	options Options
}

// http://php.net/manual/en/phar.fileformat.signature.php
func (s *signature) check(filename string) error {
	file, _ := os.Open(filename)
	defer func() {
		_ = file.Close()
	}()

	stat, err := file.Stat()
	if err != nil {
		return errors.New("can't stat file: " + err.Error())
	}

	_, err = file.Seek(-8, 2)
	if err != nil {
		return errors.New("can't seek file: " + err.Error())
	}

	type sBinary struct {
		Flag uint32
		Gbmb uint32
	}
	var sb sBinary

	err = binary.Read(file, binary.LittleEndian, &sb)
	if err != nil {
		return errors.New("can't read signature bytes: " + err.Error())
	}

	if sb.Gbmb != 1112359495 { //GBMB string
		return errors.New("can't find GBMB constant at the end")
	}

	hasher, algorithm, signatureLength := s.getHash(sb.Flag)
	if hasher == nil {
		return nil
	}

	_, err = file.Seek(-(8 + signatureLength), 2)
	if err != nil {
		return errors.New("can't seek file: " + err.Error())
	}

	fileSig := make([]byte, signatureLength)
	_, err = file.Read(fileSig)
	if err != nil {
		return errors.New("can't read file signature: " + err.Error())
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return errors.New("can't seek file: " + err.Error())
	}

	if _, err := io.CopyN(hasher, file, stat.Size()-signatureLength-8); err != nil {
		return errors.New("can't copy buffer to hasher: " + err.Error())
	}

	if !bytes.Equal(hasher.Sum(nil), fileSig) {
		return errors.New(algorithm + " hash of file is incorrect")
	}

	return nil
}

func (s *signature) getHash(flag uint32) (hash.Hash, string, int64) {
	const (
		sigMD5    = 0x0001 //PHAR_SIG_MD5
		sigSHA1   = 0x0002 //PHAR_SIG_SHA1
		sigSHA256 = 0x0003 //PHAR_SIG_SHA256
		sigSHA512 = 0x0004 //PHAR_SIG_SHA512
	)

	switch flag {
	case sigMD5:
		return md5.New(), "MD5", 16
	case sigSHA1:
		return sha1.New(), "SHA1", 20
	case sigSHA256:
		return sha256.New(), "SHA256", 32
	case sigSHA512:
		return sha512.New(), "SHA512", 64
	}

	return nil, "", 0
}
