package connector

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

type HttpConnector struct {
}

var _ Connector = (*HttpConnector)(nil)

func (h *HttpConnector) SendFile(ip string, port string, partialPath string, file io.Reader, size uint64) error {
	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", partialPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", ip+"/file/send", bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return errors.New("server returned not OK")
	}
	return nil
}
func (h *HttpConnector) RequestFile(ip string, port string, partialPath string, output io.Writer) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", ip+"/file/request", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("file", partialPath)
	req.URL.RawQuery = q.Encode()
	response, err := client.Do(req)
	if response.StatusCode != http.StatusOK {
		return errors.New("server returned not OK")
	}

	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(output, response.Body)
	return err
}
func (h *HttpConnector) RequestListFiles(ip string, port string) ([]string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", ip+"/list", nil)
	if err != nil {
		return nil, nil
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("server returned not OK HTTP status")
	}

	defer rsp.Body.Close()
	var paths []string
	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(rsp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer.Bytes(), &paths)
	if err != nil {
		return nil, err
	}

	return paths, nil
}
func (h *HttpConnector) RequestUsedSize(ip string, port string) (int64, error) {
	var sizeStr int64
	client := http.Client{}
	req, err := http.NewRequest("GET", ip+"/size", nil)
	if err != nil {
		return 0, nil
	}
	rsp, err := client.Do(req)
	if err != nil {
		return 0, nil
	}

	if rsp.StatusCode != http.StatusOK {
		return 0, errors.New("server returned not OK HTTP status")
	}

	defer rsp.Body.Close()
	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(rsp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(buffer.Bytes(), &sizeStr)
	if err != nil {
		return 0, err
	}
	return sizeStr, nil
}

func (h *HttpConnector) RemoveFile(ip string, port string, partialPath string) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", ip+"/file/remove", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("file", partialPath)
	req.URL.RawQuery = q.Encode()
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("server returned not OK HTTP status")
	}
	return nil
}
