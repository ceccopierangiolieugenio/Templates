package api

import (
	"log"
	"fmt"
	"strconv"
	"context"
	"net/http"
	"encoding/json"

	"test.90poe/services/proto_pds"
)

type WebServer struct{
	Pds PDS
}

func RunWebServer(port string, pds PDS) {
	ws := WebServer{Pds:pds}
	http.HandleFunc("/get",  ws.handlerGet)
	http.HandleFunc("/list", ws.handlerList)
	http.HandleFunc("/ping", ws.handlerPing)
	log.Printf("Staring WebServer on %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

/* Request: <URI>/get?id=ABCDE */
 func (ws *WebServer) handlerGet(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]
	if !ok {
		http.Error(w, "id field is missing!!!", http.StatusBadRequest) // 400
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	port, err := ws.Pds.Get(ctx, &proto_pds.ID{ID:id[0]})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError) // 500
		return
	}
	jsonBytes, _ := json.Marshal(port)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonBytes))
}

/* Request: <URI>/list?limit=XXX&offset=YYY
   Default if unspecified;
    limit=10
    offset=0
 */
func (ws *WebServer) handlerList(w http.ResponseWriter, r *http.Request) {
	var limit int32 = 10
	if v, ok := r.URL.Query()["limit"]; ok {
		if sv, err := strconv.Atoi(v[0]); err == nil {
			limit = int32(sv)
		}
	}
	var offset int32 = 0
	if v, ok := r.URL.Query()["offset"]; ok {
		if sv, err := strconv.Atoi(v[0]); err == nil {
			offset = int32(sv)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	port, err := ws.Pds.List(ctx, &proto_pds.Query{Limit:limit,Offset:offset})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError) // 500
		return
	}
	jsonBytes, _ := json.Marshal(port)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonBytes))
}

func (ws *WebServer) handlerPing(w http.ResponseWriter, r *http.Request) {
}