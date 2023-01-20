package buffer_test

import (
	"bytes"
	"os"
	"server/buffer"
	"testing"
)

func TestSaveRead(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		t.Error(err.Error())
	}

	testPhrase := []byte("Hello tests!")

	cacheBuffer := buffer.CreateHostFSBuffer(path + string(os.PathSeparator) + "tmp")
	defer cacheBuffer.Clean()
	testReader := bytes.NewReader(testPhrase)

	id, err := cacheBuffer.Save(testReader)

	if err != nil {
		t.Fatal(err.Error())
	}

	reader, err := cacheBuffer.Read(id)

	if err != nil {
		t.Fatal(err.Error())
	}

	finalPhrase := make([]byte, len(testPhrase))
	_, err = reader.Read(finalPhrase)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(testPhrase, finalPhrase) {
		t.Fatal("Buffer return different file after saving!")
	}
}

func TestReplace(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		t.Error(err.Error())
	}
	testPhrase := []byte("Hello tests!")
	secondPhrase := []byte("Was that replaced?")

	cacheBuffer := buffer.CreateHostFSBuffer(path + string(os.PathSeparator) + "tmp")
	defer cacheBuffer.Clean()
	testReader := bytes.NewReader(testPhrase)

	id, err := cacheBuffer.Save(testReader)

	if err != nil {
		t.Fatal(err.Error())
	}

	testReader = bytes.NewReader(secondPhrase)

	err = cacheBuffer.Replace(testReader, id)

	if err != nil {
		t.Fatal(err.Error())
	}

	reader, err := cacheBuffer.Read(id)

	if err != nil {
		t.Fatal(err.Error())
	}

	finalPhrase := make([]byte, len(secondPhrase))
	_, err = reader.Read(finalPhrase)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(secondPhrase, finalPhrase) {
		t.Fatalf("Buffer return different file after saving! Must be %s, but %s", secondPhrase, finalPhrase)
	}
}
