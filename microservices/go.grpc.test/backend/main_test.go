package main

import (
	"context"
	"testing"
	pb "testprot.eugenio/testweb"

)

/* ref: https://stackoverflow.com/questions/42102496/testing-a-grpc-service */
func BackendTest(t *testing.T) {
    s := server{}

    // set up test cases
    tests := []struct{
        name string
        want string
    } {
        {
            name: "world",
            want: "Test world",
        },
        {
            name: "123",
            want: "Test 123",
        },
    }

    for _, tt := range tests {
        req := &pb.TestRequest{Name: tt.name}
        resp, err := s.SayTest(context.Background(), req)
        if err != nil {
            t.Errorf("Test(%v) got unexpected error", err)
        }
        if resp.Message != tt.want {
            t.Errorf("Text(%v)=%v, wanted %v", tt.name, resp.Message, tt.want)
        }
    }
}