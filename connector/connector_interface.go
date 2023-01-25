package connector

import (
	"io"
)

func New(connectorType NetConnectorType) Connector {
	switch connectorType {
	case Http:
		return &HttpConnector{}
	case Tcp:
		return &TcpConnector{}
	default:
		return nil
	}
}

type Connector interface {
	SendFile(ip string, port string, partialPath string, file io.Reader, size uint64) error
	RequestFile(ip string, port string, partialPath string, output io.Writer) error
	RequestListFiles(ip string, port string) ([]string, error)
	RequestUsedSize(ip string, port string) (int64, error)
	RemoveFile(ip string, port string, partialPath string) error
}
