package netClients

import guid "github.com/google/uuid"

type netInterface interface{
	sendFile(guid.UUID) error
	requestFile(guid.UUID, string) (string, error)
	requestListFiles(guid.UUID) []string
	requestUsedSize() (int64, error)
}