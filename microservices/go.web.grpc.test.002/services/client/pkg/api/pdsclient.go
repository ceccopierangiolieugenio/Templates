package api

import (
	"os"
	"log"
	"fmt"
	"context"

	"google.golang.org/grpc"
	
	"test.90poe/services/proto_pds"
	"test.90poe/services/client/pkg/reader"
)

type PDS interface{
	Close()
	Insert(ctx context.Context, p *proto_pds.Port) (*proto_pds.ID, error)
	Get(ctx context.Context, id *proto_pds.ID) (*proto_pds.Port, error)
	List(ctx context.Context, q *proto_pds.Query) (*proto_pds.IDs, error)
	ImportFile(filename string)
}

type PDSClient struct {
	Pds proto_pds.PDSClient
	Conn *grpc.ClientConn
}

func NewPDSClient(pds_host string) (*PDSClient, error) {
	conn, err := grpc.Dial(pds_host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := proto_pds.NewPDSClient(conn)

	return &PDSClient{Pds: client, Conn: conn}, nil
}

func (pc *PDSClient) Close() {
	pc.Conn.Close()
}

func (pc *PDSClient) Insert(ctx context.Context, p *proto_pds.Port) (*proto_pds.ID, error) {
	return pc.Pds.Insert(ctx, p)
}

func (pc *PDSClient) Get(ctx context.Context, id *proto_pds.ID) (*proto_pds.Port, error) {
	return pc.Pds.Get(ctx, id)
}

func (pc *PDSClient) List(ctx context.Context, q *proto_pds.Query) (*proto_pds.IDs, error) {
	return pc.Pds.List(ctx, q)
}


func (pc *PDSClient) ImportFile(filename string) {
	pb, err := reader.NewPortBuffer(filename)
	if err == nil {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer pb.Close()
		
		for pb.More() {
			p, err := pb.Decode()
			if err != nil {
				log.Fatal(err)
			}
			rq, err := pc.Insert(ctx, p)
			if err != nil {
				fmt.Fprintf(os.Stderr,"Failed to insert: %s, %v", p.ID,err)
			}
			if p.ID != rq.GetID(){
				fmt.Fprintf(os.Stderr,"Failed to insert: %s", p.ID)
			}
		}
	}
}