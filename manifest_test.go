package phargo

import (
	"bytes"
	"testing"
)

func TestParseTruncatedHeaderShouldFail(t *testing.T) {
	m := &manifest{}
	b := bytes.NewReader([]byte{10, 20, 30})

	err := m.parse(b)
	if err == nil {
		t.Error("expecting error")
	}
}

func TestParseShouldNotExceedMaxAlasLength(t *testing.T) {
	m := &manifest{
		options: Options{
			MaxAliasLength: 10,
		},
	}
	b := bytes.NewReader([]byte{
		1, 0, 0, 0,
		0, 0, 0, 0,
		0, 1,
		0, 0, 0, 0,
		20, 0, 0, 0,
	})

	err := m.parse(b)
	if err == nil {
		t.Error("expecting error")
	}
}

func TestParseFailsReadAlias(t *testing.T) {
	m := &manifest{
		options: Options{
			MaxAliasLength: 10,
		},
	}
	b := bytes.NewReader([]byte{
		1, 0, 0, 0,
		0, 0, 0, 0,
		0, 1,
		0, 0, 0, 0,
		4, 0, 0, 0,
		65, 65, 65,
	})

	err := m.parse(b)
	if err == nil {
		t.Error("expecting error")
	}
}

func TestParseFailsReadMetadataLength(t *testing.T) {
	m := &manifest{
		options: Options{
			MaxAliasLength: 10,
		},
	}
	b := bytes.NewReader([]byte{
		1, 0, 0, 0,
		0, 0, 0, 0,
		0, 1,
		0, 0, 0, 0,
		4, 0, 0, 0,
		65, 65, 65, 65,
		4, 0, 0,
	})

	err := m.parse(b)
	if err == nil {
		t.Error("expecting error")
	}
}

func TestParseShouldNotExceedMaxMetadataLength(t *testing.T) {
	m := &manifest{
		options: Options{
			MaxAliasLength:    10,
			MaxMetaDataLength: 10,
		},
	}
	b := bytes.NewReader([]byte{
		1, 0, 0, 0,
		0, 0, 0, 0,
		0, 1,
		0, 0, 0, 0,
		4, 0, 0, 0,
		65, 65, 65, 65,
		255, 0, 0, 0,
	})

	err := m.parse(b)
	if err == nil {
		t.Error("expecting error")
	}
}

func TestParseFailsReadMetadata(t *testing.T) {
	m := &manifest{
		options: Options{
			MaxAliasLength:    10,
			MaxMetaDataLength: 10,
		},
	}
	b := bytes.NewReader([]byte{
		1, 0, 0, 0,
		0, 0, 0, 0,
		0, 1,
		0, 0, 0, 0,
		4, 0, 0, 0,
		65, 65, 65, 65,
		8, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
	})

	err := m.parse(b)
	if err == nil {
		t.Error("expecting error")
	}
}

func TestParseSuccess(t *testing.T) {
	m := &manifest{
		options: Options{
			MaxAliasLength:    10,
			MaxMetaDataLength: 10,
		},
	}
	b := bytes.NewReader([]byte{
		1, 0, 0, 0,
		2, 0, 0, 0,
		0, 1,
		0, 0, 0, 0,
		4, 0, 0, 0,
		65, 65, 65, 65,
		8, 0, 0, 0,
		66, 66, 66, 66, 66, 66, 66, 66,
	})

	err := m.parse(b)
	if err != nil {
		t.Error("unexpected error" + err.Error())
	}

	if m.EntitiesCount != 2 {
		t.Error("EntitiesCount should be 2")
	}

	if m.Version != "0.0.1" {
		t.Error("Version should be 0.0.1")
	}

	if string(m.Alias) != "AAAA" {
		t.Error("Alias should be AAAA")
	}

	if string(m.MetaSerialized) != "BBBBBBBB" {
		t.Error("MetaSerialized should be BBBBBBBB")
	}
}
