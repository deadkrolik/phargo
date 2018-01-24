package phargo

import (
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
