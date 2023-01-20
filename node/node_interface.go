package node

import (
	"server/connector"
)

type INode interface {
	GetIp() string
	GetPort() string
	GetMaxSize() uint64
	GetConnectorType() connector.NetConnectorType
}
