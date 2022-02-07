package session

import (
	"encoding/json"
	"io"
	"log"
)

type ConnectRequest struct {
	Stub string
}

type ConnectResponse struct {
	Code int
}

func writeConnectRequest(w io.Writer, stub string) error {
	log.Println("write connect request, stub", stub)
	b, err := json.Marshal(&ConnectRequest{
		Stub: stub,
	})
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func readConnectRequest(r io.Reader) (string, error) {
	buf := make([]byte, 256)
	l, err := r.Read(buf)
	if err != nil {
		return "", err
	}

	var request ConnectRequest
	err = json.Unmarshal(buf[:l], &request)
	if err != nil {
		return "", err
	}

	log.Println("read connect request, stub", request.Stub)
	return request.Stub, nil
}

func writeConnectResponse(w io.Writer, code int) error {
	log.Println("write connect response, code", code)
	b, err := json.Marshal(&ConnectResponse{
		Code: code,
	})
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func readConnectResponse(r io.Reader) (int, error) {
	buf := make([]byte, 256)
	l, err := r.Read(buf)
	if err != nil {
		return 0, err
	}

	var response ConnectResponse
	err = json.Unmarshal(buf[:l], &response)
	if err != nil {
		return 0, err
	}

	log.Println("read connect response, code", response.Code)
	return response.Code, nil
}
