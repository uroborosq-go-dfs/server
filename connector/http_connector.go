package connector

import (
	"io"
)

type HttpConnector struct {
}

var _ Connector = (*HttpConnector)(nil)

func (h *HttpConnector) SendFile(ip string, port string, partialPath string, file io.Reader, size uint64) error {
	return nil
}
func (h *HttpConnector) RequestFile(ip string, port string, partialPath string, output io.Writer) error {
	return nil
}
func (h *HttpConnector) RequestListFiles(ip string, port string) ([]string, error) {
	return nil, nil
}
func (h *HttpConnector) RequestUsedSize(ip string, port string) (uint64, error) {
	return 0, nil
}
