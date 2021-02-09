/* 
 * The MIT License (MIT)
 * 
 * Copyright 2020 Eugenio Parodi <eugenio.parodi.78@gmail.com>.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"log"
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"
	"net/http"
)

type Event struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool  // map[fieldId]true 
	FormCompletionTime int  // Seconds
}

type Dimension struct {
	Width  string
	Height string
}

type JsonPost struct {
	EventType  string `json:"eventType"`
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Pasted     bool   `json:"pasted"`
	FormId     string `json:"formId"`
	TimeTaken  int    `json:"timeTaken"`
}

type Server struct {
	Events    map[string]*Event
	ClientDir string
}

func PrettyPrint(v *Event, rw http.ResponseWriter) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Printf(string(b))
		fmt.Fprintf(rw, string(b))
	}
	return err
}

func (s *Server) CallbackHandlerGET(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(rw, req, s.ClientDir+"/index.html")
	}else{
		http.ServeFile(rw, req, s.ClientDir+req.URL.Path)
	}
}

func (s *Server) CallbackHandlerPOST(rw http.ResponseWriter, req *http.Request) {
	/* Check for the "json" content type */
	contentType := strings.Split(req.Header.Get("Content-Type"), ";")[0]
	if contentType != "application/json" {
		log.Printf("Content Type is not application/json: %s", contentType)
		http.Error(rw, "Content Type is not application/json", http.StatusUnsupportedMediaType) // 415
		return
	}

	/* Parse the json Body */
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
		http.Error(rw, "", http.StatusInternalServerError) // 500
		return
	}

	var message JsonPost
	if err = json.Unmarshal(body, &message); err != nil {
		log.Printf("Body parse error, %v", err)
		http.Error(rw, "", http.StatusBadRequest) // 400
		return
	}

	/* Check the Message integrity */
	if message.EventType=="" || message.WebsiteUrl=="" || message.SessionId=="" {
		log.Printf("Bad Request: %v", message)
		http.Error(rw, "", http.StatusBadRequest) // 400
		return
	}

	eventType  := message.EventType
	websiteUrl := message.WebsiteUrl
	sessionId  := message.SessionId

	event := s.Events[sessionId]
	if event == nil {
		/* New SessionId */
		event = &Event{
			WebsiteUrl:   websiteUrl,
			SessionId:    sessionId,
			CopyAndPaste: make(map[string]bool),
			FormCompletionTime: -1,
		}
		s.Events[sessionId] = event
	}

	switch eventType {
	case "copyAndPaste":
		if message.FormId=="" {
			log.Printf("Bad Request: %v", message)
			http.Error(rw, "Missing FormId", http.StatusBadRequest) // 400
			return
		}
		event.CopyAndPaste[message.FormId] = message.Pasted
	case "timeTaken":
		/* Negative time is the result of no character typed, 
		   i.e. only copy/paste or mouse interaction */
		if message.TimeTaken >= 0 {
			event.FormCompletionTime = message.TimeTaken/1000
		}
	case "windowSize":
		if event.ResizeFrom.Width == "" && event.ResizeFrom.Height == "" {
			event.ResizeFrom.Width  = strconv.Itoa(message.Width)
			event.ResizeFrom.Height = strconv.Itoa(message.Height)
		}else{
			event.ResizeTo.Width  = strconv.Itoa(message.Width)
			event.ResizeTo.Height = strconv.Itoa(message.Height)
		}
	default:
		log.Printf("Unrecognised eventType: %s", eventType)
		http.Error(rw, "Unrecognised eventType", http.StatusBadRequest) // 400
		return
	}

	if err = PrettyPrint(event, rw); err != nil {
		log.Printf("Print error, %v", err)
		http.Error(rw, "", http.StatusInternalServerError) // 500
		return
	}
}


func (s *Server) CallbackHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		s.CallbackHandlerGET(rw, req)
	case "POST":
		s.CallbackHandlerPOST(rw, req)
	default:
		log.Printf("Method not allowed: %s",req.Method)
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed) // 405
	}
}


func PrintHelp() {
	fmt.Println("HTTP Server :")
	flag.PrintDefaults()
}

func getFlags(s *Server, port *string) bool{
	var printHelp bool
	flag.BoolVar(&printHelp, "help", false, "Print this help message.")
	flag.StringVar(&s.ClientDir, "clientdir", "./client", "Client Directory")
	flag.StringVar(port, "port", "5000", "http port")
	flag.Parse()
	return printHelp
}

func main() {
	var port string
	srv := &Server{Events:make(map[string]*Event),}
	if getFlags(srv, &port){
		PrintHelp()
		return
	}

	http.HandleFunc("/", srv.CallbackHandler)
	log.Printf("Staring WebServer on %s...",port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}