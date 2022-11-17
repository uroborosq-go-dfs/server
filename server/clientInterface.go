package server

import guid "github.com/google/uuid"


type clientInterface interface {
	addFile(path string) error
	removeFile(partialPath string) error
	addNode(ip string, port string, size uint64) (guid.UUID, error)
	removeNode(guid.UUID) error
	cleanNode(guid.UUID) error
	listOfNodes() ([]guid.UUID, []int64, error)
	listOfFiles(guid.UUID) error
	balanceNodes() error
}