package server

import


func createServer(tmpPath string, ) {
	
}


type server struct {
	
}

func (s server) addFile(path string) error {
	return nil
}

func (s server) removeFile(partialPath string) error {
	return nil
}

func (s server) addNode(ip string, port string, size uint64) (guid.UUID, error) {
	
}

func (s server) removeNode(guid.UUID) error {
	return nil
}

func (s server) cleanNode(guid.UUID) error {

}

func (s server) listOfNodes() ([]guid.UUID, []int64, error) {

}

func (s server) listOfFiles(guid.UUID) error {

}

func (s server) balanceNodes() error {

}