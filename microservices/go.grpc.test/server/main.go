/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"fmt"
	"time"
    "net/http"

	"google.golang.org/grpc"
	pb "testprot.eugenio/testweb"
)

var address string = "localhost:50051"

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ln := r.URL.Path[1:]
	log.Printf("Handling: %s", ln)    

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := ln

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rq, err := c.SayTest(ctx, &pb.TestRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", rq.GetMessage())
	
	fmt.Fprintf(w, "Greeting: %s", rq.GetMessage())
}

func main() {
	if (len(os.Args[1:])==1){
		/* I assume a valid the backend uri is used as arg */
		address = os.Args[1]
	}

	http.HandleFunc("/", CallbackHandler)
	log.Println("Using the backend at:", address)
	log.Println("Staring WebServer on 5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}