package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	//"io"
	"log"
	"os"

	"test.90poe/services/proto_pds"
)

func PortScanner() {

}

type PortBuffer struct{
	File   *os.File
	Reader *bufio.Reader
	Dec    *json.Decoder
	NextP  *proto_pds.Port
}

func NewPortBuffer(filename string) (*PortBuffer, error){
	pb := &PortBuffer{}

	var err error
	pb.File, err = os.Open(filename)
	if err != nil {
		log.Printf("error reading from file %s: %V", filename, err)
		return nil, err
	}

	pb.Reader = bufio.NewReader(pb.File)
	pb.Dec = json.NewDecoder(pb.Reader)

	// read open bracket
	t, err := pb.Dec.Token()
	if err != nil || fmt.Sprint(t) != "{" {
		return nil, err
	}

	return pb, nil
}

func (pb *PortBuffer) Close(){
	if pb.File != nil {
		pb.File.Close()
	}
}

func (pb *PortBuffer) More() bool {
	// read ID or Close bracket
	t, err := pb.Dec.Token()
	id := fmt.Sprint(t)
	if err != nil || id == "}" {
		return false
	}
	var content proto_pds.Content
	err = pb.Dec.Decode(&content)
	if err != nil {
		fmt.Fprintln(os.Stderr,"error decoding to content: %s", err)
		return false
	}
	pb.NextP = &proto_pds.Port{ID:id, Content:&content}
	return true
}

func (pb *PortBuffer) Decode() (*proto_pds.Port, error){
	ret := pb.NextP
	pb.NextP = nil
	return ret, nil
}


func NewPortDecoder() {

}
