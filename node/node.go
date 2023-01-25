package node

import (
	"errors"

	"github.com/uroborosq-go-dfs/server/connector"
)

func CreateNode(ip string, port string, maxSize int64, connectType connector.NetConnectorType) INode {
	return &Node{ip, port, 0, maxSize, connectType}
}

type Node struct {
	ip          string
	port        string
	currentSize int64
	maximumSize int64
	connectType connector.NetConnectorType
}

var _ INode = (*Node)(nil)

func (n *Node) GetIp() string {
	return n.ip
}

func (n *Node) GetPort() string {
	return n.port
}

func (n *Node) GetCurrentSize() int64 {
	return n.currentSize
}

func (n *Node) UpdateCurrentSize(add int64) error {
	if add < 0 {
		return errors.New("size must be greater than zero")
	}
	n.currentSize = add
	return nil
}

func (n *Node) GetMaxSize() int64 {
	return n.maximumSize
}

func (n *Node) GetConnectorType() connector.NetConnectorType {
	return n.connectType
}
