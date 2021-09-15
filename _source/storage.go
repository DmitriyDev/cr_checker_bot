package main

import (
	"io/ioutil"
	"log"
	"os"
)

type Storage struct {
	path string
}

func (s Storage) New(path string) Storage {
	st := Storage{path: path}
	if !st.exists() {
		st.create()
	}

	return st
}

func (s *Storage) exists() bool {
	_, err := os.Stat(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func (s *Storage) create() {
	newFile, err := os.Create(s.path)
	if err != nil {
		log.Fatal(err)
	}

	err = newFile.Chmod(0755)
	if err != nil {
		log.Fatal(err)
	}

	newFile.Close()
}

func (s *Storage) fileResource(flag int, mode os.FileMode) *os.File {
	f, err := os.OpenFile(s.path, flag, mode)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func (s *Storage) Get() string {
	f := s.fileResource(os.O_RDONLY, os.ModePerm)
	defer f.Close()
	bData, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Get Error: ", err)
	}

	return string(bData)
}

func (s *Storage) Append(data string) {
	f := s.fileResource(os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()

	_, err := f.WriteString(data + "\n")

	if err != nil {
		log.Fatal("Append Error: ", err)
	}

}

func (s *Storage) ReplaceWith(data string) {
	s.create()
	f := s.fileResource(os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	_, err := f.WriteString(data)

	if err != nil {
		log.Fatal(err)
	}
}
