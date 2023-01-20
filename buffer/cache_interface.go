package buffer

import (
	guid "github.com/google/uuid"

	"io"
)

type Buffer interface {
	Save(io.Reader) (guid.UUID, error)
	Replace(io.Reader, guid.UUID) error
	Read(guid.UUID) (io.Reader, error)
	Clean() error
}
