package api

import (
	"context"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"


	"test.90poe/services/proto_pds"
)

type PDSMock struct {
	Port *proto_pds.Port
	IDs *proto_pds.IDs
}

func (pc *PDSMock) Close() {
}

func (pm *PDSMock) Insert(ctx context.Context, p *proto_pds.Port) (*proto_pds.ID, error) {
	return &proto_pds.ID{ID:pm.Port.ID}, nil
}

func (pm *PDSMock) Get(ctx context.Context, id *proto_pds.ID) (*proto_pds.Port, error) {
	return pm.Port, nil
}

func (pm *PDSMock) List(ctx context.Context, q *proto_pds.Query) (*proto_pds.IDs, error) {
	return &proto_pds.IDs{}, nil
}

func (pm *PDSMock) ImportFile(filename string) {
}

/* ref: https://blog.questionable.services/article/testing-http-handlers-go/ */
func TestGet(t *testing.T) {
	port := &proto_pds.Port{
		ID: "ABCDE",
		Content: &proto_pds.Content{
				Name: "Test ABCDE",
				City: "Test City",
				Country: "Test Country",
			},
	}
	ws := &WebServer{Pds:&PDSMock{Port: port}}
	
	req, err := http.NewRequest("GET", "/get?id=ABCDE", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.handlerGet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	jsonBytes, _ := json.Marshal(port)
	if rr.Body.String() != string(jsonBytes) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(jsonBytes))
	}
}