package main

import (
	"bytes"
	"testing"
)

func TestReadWriterComposition(t *testing.T) {
	var buf bytes.Buffer
	var rw ReadWriter = &buf
	rw.Write([]byte("test"))
	readBuf := make([]byte, 4)
	n, err := rw.Read(readBuf)
	if err != nil {
		t.Fatal(err)
	}
	if string(readBuf[:n]) != "test" {
		t.Errorf("expected 'test', got '%s'", string(readBuf[:n]))
	}
}

func TestUserValidatedFormatter(t *testing.T) {
	var vf ValidatedFormatter = User{Name: "Alice", Email: "a@b.com"}
	if err := vf.Validate(); err != nil {
		t.Fatal(err)
	}
	formatted := vf.Format()
	if formatted != "Alice <a@b.com>" {
		t.Errorf("expected 'Alice <a@b.com>', got '%s'", formatted)
	}
}

func TestUserValidateFailsOnEmptyName(t *testing.T) {
	u := User{Name: "", Email: "a@b.com"}
	err := u.Validate()
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUserValidateFailsOnEmptyEmail(t *testing.T) {
	u := User{Name: "Alice", Email: ""}
	err := u.Validate()
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestFileStorageImplementsStorage(t *testing.T) {
	var s Storage = &FileStorage{}
	_ = s
}

func TestBackup(t *testing.T) {
	fs := &FileStorage{}
	Backup(fs)
	if fs.data != "loaded data" {
		t.Errorf("expected 'loaded data', got '%s'", fs.data)
	}
}

func TestValidatorInterface(t *testing.T) {
	var v Validator = User{Name: "Alice", Email: "a@b.com"}
	if err := v.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestFormatterInterface(t *testing.T) {
	var f Formatter = User{Name: "Alice", Email: "a@b.com"}
	result := f.Format()
	if result != "Alice <a@b.com>" {
		t.Errorf("expected 'Alice <a@b.com>', got '%s'", result)
	}
}

func TestUserValidateMultipleErrors(t *testing.T) {
	u := User{Name: "", Email: ""}
	err := u.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
}
