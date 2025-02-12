package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	ID       int
	Username string
}

type Server struct {
	db    map[int]*User
	dbHit int
}

func NewServer() *Server {
	db := make(map[int]*User)

	for i := 0; i < 100; i++ {
		db[i] = &User{
			ID:       i + 1,
			Username: fmt.Sprintf("user_%d", i+1),
		}

	}
	return &Server{
		db: db,
	}
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//hit the database
	user, ok := s.db[id]
	if !ok {
		panic("user not found")
	}
	s.dbHit++
	json.NewEncoder(w).Encode(user)
}
