package server

import (
	"errors"
	guid "github.com/google/uuid"
	"os"
	"server/connector"
	"server/node"
)

const MaxInt64 = ^uint64(0)

func CreateServer(tmpPath string) *Server {
	return &Server{
		nodes: make(map[guid.UUID]node.INode),
		sizes: make(map[guid.UUID]uint64),
	}
}

type Server struct {
	nodes map[guid.UUID]node.INode
	sizes map[guid.UUID]uint64
}

func (s *Server) AddFile(path string, nodePath string) error {
	if len(s.sizes) == 0 {
		return errors.New("firstly add nodes to store files")
	}
	min := MaxInt64
	id := guid.Nil
	for i, size := range s.sizes {
		if size < min {
			id = i
			min = size
		}
	}

	connect := connector.CreateConnector(s.nodes[id].GetConnectorType())

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	return connect.SendFile(s.nodes[id].GetIp(), s.nodes[id].GetPort(), nodePath, file, uint64(info.Size()))
}

func (s *Server) RemoveFile(partialPath string) error {
	for _, n := range s.nodes {
		connect := connector.CreateConnector(n.GetConnectorType())
		paths, err := connect.RequestListFiles(n.GetIp(), n.GetPort())
		if err != nil {
			return err
		}
		for _, path := range paths {
			if path == partialPath {
				//TODO: implement deleting
			}
		}
	}
	return errors.New("not implemented")
}

func (s *Server) AddNode(ip string, port string, size uint64, connectType connector.NetConnectorType) (guid.UUID, error) {
	n := node.CreateNode(ip, port, size, connectType)
	newGuid := guid.New()

	connect := connector.CreateConnector(connectType)
	usedSize, err := connect.RequestUsedSize(ip, port)
	if err != nil {
		return guid.Nil, err
	}
	s.sizes[newGuid] = usedSize
	s.nodes[newGuid] = n
	return newGuid, nil
}

func (s *Server) RemoveNode(id guid.UUID) {
	delete(s.nodes, id)
	delete(s.sizes, id)
}

func (s *Server) CleanNode(guid.UUID) error {
	return nil
}

func (s *Server) ListOfAllFiles() ([]string, error) {
	paths := make([]string, 0)
	for _, n := range s.nodes {
		connect := connector.CreateConnector(n.GetConnectorType())
		nodePaths, err := connect.RequestListFiles(n.GetIp(), n.GetPort())
		if err != nil {
			return nil, err
		}
		paths = append(paths, nodePaths...)
	}
	return paths, nil
}

func (s *Server) ListOfNodeFiles(id guid.UUID) ([]string, error) {
	n, ok := s.nodes[id]
	if !ok {
		return nil, errors.New("no such node")
	}
	connect := connector.CreateConnector(n.GetConnectorType())
	return connect.RequestListFiles(n.GetIp(), n.GetPort())
}

func (s *Server) BalanceNodes() error {
	return nil
}
