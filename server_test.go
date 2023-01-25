package main_test

import (
	"fmt"
	"github.com/uroborosq-go-dfs/server/server"
	"os"
	"path/filepath"
	"testing"
)

var (
	host     = os.Getenv("GODFS_DB_HOST")
	port     = os.Getenv("GODFS_DB_PORT")
	user     = os.Getenv("GODFS_DB_USER")
	password = os.Getenv("GODFS_DB_PASSWORD")
	dbname   = os.Getenv("GODFS_DB_DBNAME")
	driver   = os.Getenv("GODFS_DB_DRIVER")
)

func TestBenchTest(t *testing.T) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	s, err := server.CreateServer(driver, conn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = filepath.Walk(os.Getenv("GODFS_TEST_FILEPATH"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return s.AddFile(path, info.Name())
		} else {
			return nil
		}
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func TestSendSingleFile(t *testing.T) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	s, err := server.CreateServer(driver, conn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	_, err = s.AddNode(os.Getenv("GODFS_TEST_NODE_IP"), os.Getenv("GODFS_TEST_NODE_PORT"), 11, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = s.AddFile(os.Getenv("GODFS_TEST_FILE"), "test")
	if err != nil {
		t.Fatal(err.Error())
	}
}
