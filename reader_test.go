package phargo

import (
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	r := NewReader()

	file, err := r.Parse("./testdata/simple.phar")
	if err != nil {
		t.Error("Got error", err)
		return
	}

	if len(file.Files) != 2 {
		t.Error("Not 2 files")
		return
	}

	if file.Files[0].Name != "1.txt" || string(file.Files[0].Data) != "ASDF" {
		t.Error("Wrong 1 file content or name")
		return
	}

	if file.Files[1].Name != "index.php" || string(file.Files[1].Data) != "ZXCV" {
		t.Error("Wrong 2 file content or name")
		return
	}

	if string(file.Metadata) != "a:1:{s:1:\"a\";i:123;}" {
		t.Error("Wrong metadata")
		return
	}
}

func TestAliasMD5(t *testing.T) {
	r := NewReader()

	file, err := r.Parse("./testdata/alias_md5.phar")
	if err != nil {
		t.Error("Got error", err)
		return
	}

	if len(file.Files) != 1 {
		t.Error("Not 1 file")
		return
	}

	if file.Files[0].Name != "data.txt" || string(file.Files[0].Data) != "DATA" {
		t.Error("Wrong 1 file content or name")
		return
	}

	if file.Alias != "ALIAS" {
		t.Error("Incorrect alias")
		return
	}
}

func TestMetadataDirSHA256(t *testing.T) {
	r := NewReader()

	file, err := r.Parse("./testdata/metadata_dir_sha256.phar")
	if err != nil {
		t.Error("Got error", err)
		return
	}

	if len(file.Files) != 4 {
		t.Error("Not 4 files")
		return
	}

	if file.Files[0].Name != "FILE" || string(file.Files[0].Data) != "FDATA" {
		t.Error("Wrong 1 file content or name")
		return
	}

	if file.Files[1].Name != "DIR1/FILE1" || string(file.Files[1].Data) != "D1_DATA11" {
		t.Error("Wrong 2 file content or name")
		return
	}

	if file.Files[2].Name != "DIR1/FILE2" || string(file.Files[2].Data) != "D1_DATA12" {
		t.Error("Wrong 3 file content or name")
		return
	}

	if file.Files[3].Name != "DIR2/FILE1" || string(file.Files[3].Data) != "D1_DATA21" {
		t.Error("Wrong 3 file content or name")
		return
	}

	if string(file.Files[0].Metadata) != "a:1:{s:1:\"v\";s:1:\"x\";}" {
		t.Error("Wrong metadata for file 1")
		return
	}

	if string(file.Files[1].Metadata) != "" {
		t.Error("Wrong metadata for file 2")
		return
	}

	if string(file.Files[3].Metadata) != "a:1:{s:1:\"z\";s:2:\"cc\";}" {
		t.Error("Wrong metadata for file 3")
		return
	}
}

func TestBadHash(t *testing.T) {
	r := NewReader()

	_, err := r.Parse("./testdata/bad_hash.phar")
	if err == nil {
		t.Error("Should get error")
		return
	}

	if !strings.Contains(err.Error(), "MD5 hash") {
		t.Error("Should be MD5 hash error")
		return
	}
}

func TestSHA512(t *testing.T) {
	r := NewReader()

	file, err := r.Parse("./testdata/sha512.phar")
	if err != nil {
		t.Error("Got error", err)
		return
	}

	if len(file.Files) != 1 {
		t.Error("Not 1 file")
		return
	}
}

func TestGZ(t *testing.T) {
	r := NewReader()

	file, err := r.Parse("./testdata/gz.phar")
	if err != nil {
		t.Error("Got error", err)
		return
	}

	if len(file.Files) != 1 {
		t.Error("Not 1 file")
		return
	}

	if string(file.Files[0].Data) != "DATADATADATADATA" {
		t.Error("Wrong file 1 content")
		return
	}
}
