package buffer

import (
	"errors"
	"io"
	"os"

	guid "github.com/google/uuid"
)

func CreateHostFSBuffer(path string) *HostFSBuffer {
	return &HostFSBuffer{path}
}

type HostFSBuffer struct {
	path string
}

var _ Buffer = (*HostFSBuffer)(nil)

func (h *HostFSBuffer) Save(reader io.Reader) (guid.UUID, error) {
	newGuid := guid.New()
	err := os.MkdirAll(h.path, os.ModePerm)

	if err != nil {
		return guid.Nil, err
	}

	file, err := os.OpenFile(h.path+string(os.PathSeparator)+newGuid.String(), os.O_CREATE|os.O_RDWR, 0611)
	defer file.Close()

	if err != nil {
		return guid.Nil, err
	}

	pieceOfData := [1024]byte{}
	for n, readerErr := reader.Read(pieceOfData[:]); n != 0 && readerErr == nil; n, readerErr = reader.Read(pieceOfData[:]) {
		write, err := file.Write(pieceOfData[:n])
		if err != nil {
			return guid.Nil, err
		} else if write != n {
			return guid.Nil, errors.New("written bytes don't match read bytes")
		}
	}

	return newGuid, nil
}

func (h *HostFSBuffer) Replace(reader io.Reader, id guid.UUID) error {
	file, err := os.OpenFile(h.path+string(os.PathSeparator)+id.String(), os.O_RDWR, 0610)

	if err != nil {
		return err
	}

	pieceOfData := [1024]byte{}
	for n, readerErr := reader.Read(pieceOfData[:]); n != 0 && readerErr == nil; n, readerErr = reader.Read(pieceOfData[:]) {
		write, err := file.Write(pieceOfData[:n])
		if err != nil {
			return err
		} else if write != n {
			return errors.New("written bytes don't match read bytes")
		}
	}

	return nil
}

func (h *HostFSBuffer) Read(guid guid.UUID) (io.Reader, error) {
	file, err := os.Open(h.path + string(os.PathSeparator) + guid.String())
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (h *HostFSBuffer) Clean() error {
	return os.RemoveAll(h.path)
}
