package node

import (
	"server/connector"
)

func CreateNode(ip string, port string, maxSize uint64, connectType connector.NetConnectorType) INode {
	return &Node{ip, port, maxSize, connectType}
}

type Node struct {
	ip          string
	port        string
	maximumSize uint64
	connectType connector.NetConnectorType
}

var _ INode = (*Node)(nil)

func (n *Node) GetIp() string {
	return n.ip
}

func (n *Node) GetPort() string {
	return n.port
}

func (n *Node) GetMaxSize() uint64 {
	return n.maximumSize
}

func (n *Node) GetConnectorType() connector.NetConnectorType {
	return n.connectType
}
