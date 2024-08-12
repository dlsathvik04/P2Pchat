package api

import (
	"net/http"

	"github.com/dlsathvik04/P2Pchat/db"
	"github.com/dlsathvik04/P2Pchat/pkg/hasher"
	"github.com/dlsathvik04/P2Pchat/pkg/jwt"
)

type APIServer interface {
	ListenAndServe()
}

type P2PChatAPIServer struct {
	hostname   string
	mux        *http.ServeMux
	p2pchatDB  db.P2PchatDB
	hasher     hasher.Hasher
	jwtManager jwt.JWTManager
}

func NewAPIServer(hostname string, mux *http.ServeMux, p2pdb db.P2PchatDB, hasher hasher.Hasher, jwtman jwt.JWTManager) APIServer {
	server := P2PChatAPIServer{hostname, mux, p2pdb, hasher, jwtman}
	server.setUp()
	return &server

}

func (s *P2PChatAPIServer) ListenAndServe() {
	http.ListenAndServe(s.hostname, s.mux)
}

func (s *P2PChatAPIServer) setUp() {

	authManager := NewAuthAPIManager(s.hasher, s.jwtManager, &s.p2pchatDB)
	authManager.Register(s.mux)

	stunApiManger := NewSTUNAPIManager(s.p2pchatDB, authManager)
	stunApiManger.Register(s.mux)

	rendApiManager := NewRendezvousAPIManager(s.p2pchatDB, authManager)
	rendApiManager.Register(s.mux)
}
