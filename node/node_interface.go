package node

import (
	"github.com/uroborosq-go-dfs/server/connector"
)

type INode interface {
	GetIp() string
	GetPort() string
	GetCurrentSize() int64
	UpdateCurrentSize(int64) error
	GetMaxSize() int64
	GetConnectorType() connector.NetConnectorType
}
