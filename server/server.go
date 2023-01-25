package server

import (
	"errors"
	"os"

	guid "github.com/google/uuid"
	"github.com/uroborosq-go-dfs/server/connector"
	"github.com/uroborosq-go-dfs/server/node"
	node_storage "github.com/uroborosq-go-dfs/server/node-storage"
)

const MaxInt64 = int64(9223372036854775807)

func CreateServer(dbType string, uri string) (*Server, error) {
	nodeStorage, err := node_storage.New(dbType, uri)
	if err != nil {
		return nil, err
	}
	pathStorage, err := node_storage.NewPathStorage(dbType, uri)
	if err != nil {
		return nil, err
	}
	return &Server{
		nodes: nodeStorage,
		paths: pathStorage,
	}, nil
}

type Server struct {
	nodes *node_storage.NodeStorage
	paths *node_storage.PathStorage
}

func (s *Server) AddFile(path string, nodePath string) error {
	ids, nodes, err := s.nodes.GetAll()
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return errors.New("firstly add nodes to store files")
	}
	min := MaxInt64
	index := -1
	for i, n := range nodes {
		if n.GetCurrentSize() < min {
			index = i
			min = n.GetCurrentSize()
		}
	}

	connect := connector.New(nodes[index].GetConnectorType())

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	err = s.paths.Add(ids[index], nodePath, info.Size())
	if err != nil {
		return err
	}
	err = nodes[index].UpdateCurrentSize(nodes[index].GetCurrentSize() + info.Size())
	if err != nil {
		return err
	}
	err = s.nodes.Replace(ids[index], nodes[index])
	if err != nil {
		return err
	}
	return connect.SendFile(nodes[index].GetIp(), nodes[index].GetPort(), nodePath, file, uint64(info.Size()))
}

func (s *Server) RemoveFile(partialPath string) error {
	id, size, err := s.paths.Get(partialPath)
	if err != nil {
		return err
	}
	n, err := s.nodes.Get(id)
	if err != nil {
		return err
	}
	err = n.UpdateCurrentSize(n.GetCurrentSize() - size)
	if err != nil {
		return err
	}
	connect := connector.New(n.GetConnectorType())
	err = connect.RemoveFile(n.GetIp(), n.GetPort(), partialPath)
	if err != nil {
		return err
	}
	err = s.paths.Remove(partialPath)
	if err != nil {
		return err
	}
	err = s.nodes.Replace(id, n)
	return err
}

func (s *Server) AddNode(ip string, port string, size int64, connectType connector.NetConnectorType) (guid.UUID, error) {
	n := node.CreateNode(ip, port, size, connectType)
	newGuid := guid.New()

	connect := connector.New(connectType)
	usedSize, err := connect.RequestUsedSize(ip, port)
	if err != nil {
		return guid.Nil, err
	}
	err = n.UpdateCurrentSize(usedSize)
	if err != nil {
		return guid.Nil, err
	}
	err = s.nodes.Add(newGuid, n)
	return newGuid, err
}

func (s *Server) RemoveNode(id guid.UUID) error {
	err := s.nodes.Remove(id)

	paths, _, err := s.paths.GetPathsByNodeId(id)
	for _, path := range paths {
		err = s.paths.Remove(path)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *Server) CleanNode(id guid.UUID) error {
	n, err := s.nodes.Get(id)
	if err != nil {
		return err
	}
	connect := connector.New(n.GetConnectorType())
	paths, _, err := s.paths.GetPathsByNodeId(id)
	if err != nil {
		return err
	}
	for _, path := range paths {
		err = connect.RemoveFile(n.GetIp(), n.GetPort(), path)
		if err != nil {
			return err
		}
		err = s.paths.Remove(path)
		if err != nil {
			return err
		}
	}
	_ = n.UpdateCurrentSize(0)
	err = s.nodes.Replace(id, n)
	return err
}

func (s *Server) ListOfAllFiles() ([]string, []int64, error) {
	_, paths, sizes, err := s.paths.GetAll()
	return paths, sizes, err
}

func (s *Server) ListOfNodeFiles(id guid.UUID) ([]string, []int64, error) {
	paths, sizes, err := s.paths.GetPathsByNodeId(id)
	if err != nil {
		return nil, nil, err
	}
	return paths, sizes, err
}

func (s *Server) GetFile(partialPath string, pathToDownload string) error {
	id, _, err := s.paths.Get(partialPath)
	if err != nil {
		return err
	}
	n, err := s.nodes.Get(id)
	if err != nil {
		return err
	}
	connect := connector.New(n.GetConnectorType())

	file, err := os.OpenFile(pathToDownload, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	return connect.RequestFile(n.GetIp(), n.GetPort(), partialPath, file)
}
