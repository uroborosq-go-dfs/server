package connector

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"

	models "github.com/uroborosq-go-dfs/models/tcp-operation-code"
)

type TcpConnector struct {
}

var _ Connector = (*TcpConnector)(nil)

func (h *TcpConnector) SendFile(ip string, port string, partialPath string, file io.Reader, size uint64) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	reader := bufio.NewReader(file)
	if err != nil {
		return err
	}
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	write, err := writer.Write([]byte{models.SendFile})

	if err != nil {
		return err
	} else if write != 1 {
		return errors.New("can't write to the stream")
	}

	pathLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(pathLenBytes, uint32(len(partialPath)))
	write, err = writer.Write(pathLenBytes)

	if err != nil {
		return err
	} else if write != 4 {
		return errors.New("can't write to the stream")
	}

	write, err = writer.Write([]byte(partialPath))

	if err != nil {
		return err
	} else if write != len(partialPath) {
		return errors.New("can't write to the stream")
	}

	pathLenBytes = make([]byte, 8)
	binary.BigEndian.PutUint64(pathLenBytes, size)
	write, err = writer.Write(pathLenBytes)

	if err != nil {
		return err
	} else if write != 8 {
		return errors.New("can't write to the stream")
	}

	bufferSize := uint64(1024)
	buffer := make([]byte, 1024)
	for i := uint64(0); i < size; i += bufferSize {
		read, err := reader.Read(buffer)
		if err != nil {
			return err
		}
		_, err = writer.Write(buffer[:read])

		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
func (h *TcpConnector) RequestFile(ip string, port string, partialPath string, output io.Writer) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(conn)
	write, err := writer.Write([]byte{models.RequestFile})
	if err != nil {
		return err
	} else if write != 1 {
		return errors.New("can't write to stream")
	}
	pathLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(pathLenBytes, uint32(len(partialPath)))
	write, err = writer.Write(pathLenBytes)
	if err != nil {
		return err
	} else if write != 4 {
		return errors.New("can't write to stream")
	}
	write, err = writer.Write([]byte(partialPath))
	if err != nil {
		return err
	} else if write != len(partialPath) {
		return errors.New("can't write to stream")
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	bufferSize := uint64(1024)
	buffer := make([]byte, bufferSize)
	pathLenBytes = make([]byte, 8)
	read, err := reader.Read(pathLenBytes)
	if err != nil {
		return err
	} else if read != 8 {
		return errors.New("can't read from stream")
	}
	fileSize := binary.BigEndian.Uint64(pathLenBytes)
	for i := uint64(0); i < fileSize; i += bufferSize {
		read, err = reader.Read(buffer)
		if err != nil {
			return err
		}
		write, err = output.Write(buffer[:read])
		if err != nil {
			return err
		} else if write != read {
			return errors.New("can't write to the stream")
		}
	}

	return nil
}
func (h *TcpConnector) RequestListFiles(ip string, port string) ([]string, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(conn)
	write, err := writer.Write([]byte{models.RequestList})
	if err != nil {
		return nil, err
	} else if write != 1 {
		return nil, errors.New("can't write to stream")
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	replySize := make([]byte, 8)
	read, err := reader.Read(replySize)
	if err != nil {
		return nil, err
	} else if read != 8 {
		return nil, errors.New("can't read from stream")
	}
	size := binary.BigEndian.Uint64(replySize)
	bufferSize := uint64(1024)
	buffer := make([]byte, bufferSize)
	pathList := make([]string, 0)
	s := ""
	for i := uint64(0); i < size; i += bufferSize {
		read, err = reader.Read(buffer)
		if err != nil {
			return nil, err
		}
		for j := 0; j < read; j++ {
			if buffer[j] == byte('\n') || buffer[j] == 0 {
				pathList = append(pathList, s)
			} else {
				s += string(buffer[j])
			}
		}
	}
	return pathList, nil
}
func (h *TcpConnector) RequestUsedSize(ip string, port string) (int64, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return 0, err
	}
	writer := bufio.NewWriter(conn)
	write, err := writer.Write([]byte{models.RequestSize})
	if err != nil {
		return 0, err
	} else if write != 1 {
		return 0, errors.New("can't write to the stream")
	}
	err = writer.Flush()
	if err != nil {
		return 0, err
	}
	reader := bufio.NewReader(conn)
	sizeByte := make([]byte, 8)
	read, err := reader.Read(sizeByte)
	if err != nil {
		return 0, err
	} else if read != 8 {
		return 0, errors.New("can't read from the stream")
	}
	return int64(binary.BigEndian.Uint64(sizeByte)), nil
}
func (h *TcpConnector) RemoveFile(ip string, port string, partialPath string) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(conn)
	write, err := writer.Write([]byte{models.RemoveFile})
	if err != nil {
		return err
	} else if write != 1 {
		return errors.New("can't write to stream")
	}
	pathLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(pathLenBytes, uint32(len(partialPath)))
	write, err = writer.Write(pathLenBytes)
	if err != nil {
		return err
	} else if write != 4 {
		return errors.New("can't write to stream")
	}
	write, err = writer.Write([]byte(partialPath))
	if err != nil {
		return err
	} else if write != len(partialPath) {
		return errors.New("can't write to stream")
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
