package main

import (
	"bytes"
	"errors"
	"fmt"
)

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

type Validator interface {
	Validate() error
}

type Formatter interface {
	Format() string
}

type ValidatedFormatter interface {
	Validator
	Formatter
}

type User struct {
	Name  string
	Email string
}

func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (u User) Format() string {
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

type Saver interface {
	Save() error
}

type Loader interface {
	Load() error
}

type Storage interface {
	Saver
	Loader
}

type FileStorage struct {
	data string
}

func (f *FileStorage) Save() error {
	fmt.Println("saving:", f.data)
	return nil
}

func (f *FileStorage) Load() error {
	f.data = "loaded data"
	return nil
}

func Backup(s Storage) {
	s.Save()
	s.Load()
}

func main() {
	var buf bytes.Buffer
	var rw ReadWriter = &buf
	rw.Write([]byte("hello"))
	var readBuf [16]byte
	n, _ := rw.Read(readBuf[:])
	fmt.Println("ReadWriter:", string(readBuf[:n]))

	var vf ValidatedFormatter = User{Name: "Alice", Email: "alice@example.com"}
	if err := vf.Validate(); err != nil {
		fmt.Println("validation failed:", err)
		return
	}
	fmt.Println("formatted:", vf.Format())

	fs := &FileStorage{data: "my data"}
	Backup(fs)
	fmt.Println("loaded data:", fs.data)
}
