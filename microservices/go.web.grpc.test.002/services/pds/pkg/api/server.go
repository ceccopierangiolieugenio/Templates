package api

import (
	"os"
	"fmt"
	"context"
	"errors"
	"encoding/json"

	"database/sql"
	_ "github.com/lib/pq"


	"test.90poe/services/pds/pkg/config"
	"test.90poe/services/proto_pds"
)

type Server struct {
	proto_pds.UnimplementedPDSServer
	Config *config.Config
	Db *sql.DB
}

func NewServer(cfg *config.Config) (*Server, error){
    db, err := sql.Open("postgres", "postgres://"+cfg.SQL_User+":"+cfg.SQL_Pass+"@"+cfg.SQL_Host+"/?sslmode=disable")
    if err != nil {
        return nil, err
    }

	srv := &Server{Config: cfg, Db: db}
	return srv, nil
}

func (s *Server) CloseDbConnections(){
	s.Db.Close()
}

func (s *Server) Insert(ctx context.Context, p *proto_pds.Port) (*proto_pds.ID, error) {
	if p == nil {
		return nil, errors.New("port is nil, aborting")
	}

	jsonBytes, _ := json.Marshal(p.Content)
	_, err := s.Db.Exec(
				`INSERT INTO ports90 VALUES ( $1, $2 ) `+
				`ON CONFLICT (id) DO UPDATE SET content=excluded.content ;`,
				p.ID, string(jsonBytes))

	if err != nil {
		fmt.Fprintln(os.Stderr,"SQL error: ", err)
		return nil, errors.New("SQL error")
	}
	return &proto_pds.ID{ID:p.ID}, nil
}

func (s *Server) Get(ctx context.Context, in *proto_pds.ID) (*proto_pds.Port, error) {
	var content string
	s.Db.QueryRow(`SELECT content FROM  ports90 WHERE id=$1 ;`,in.ID).Scan(&content)
	ret := &proto_pds.Port{ID:in.ID}
	json.Unmarshal([]byte(content), &ret.Content)
	return ret, nil
}

func (s *Server) List(ctx context.Context, in *proto_pds.Query) (*proto_pds.IDs, error) {
	rows, err := s.Db.Query(`SELECT id FROM  ports90 LIMIT $1 OFFSET $2 ;`,in.Limit, in.Offset)
	if err != nil {
		fmt.Fprintln(os.Stderr,"SQL error: ", err)
		return nil, errors.New("SQL error")
	}
	defer rows.Close()

	var ret proto_pds.IDs
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			fmt.Fprintln(os.Stderr,"SQL error: ", err)
			return nil, errors.New("SQL error")
		}
		ret.IDs = append(ret.IDs,&proto_pds.ID{ID:id})
	}

	return &ret, nil
}