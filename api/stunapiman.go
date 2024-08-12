package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dlsathvik04/P2Pchat/db"
)

type STUNAPIManager struct {
	store db.P2PchatDB
	aam   *AuthAPIManager
}

func NewSTUNAPIManager(store db.P2PchatDB, aam *AuthAPIManager) *STUNAPIManager {
	return &STUNAPIManager{store, aam}
}

func (s *STUNAPIManager) Register(router *http.ServeMux) {
	router.HandleFunc("GET /stun/", s.aam.Authorize(s.handleStun))
}

func (s *STUNAPIManager) handleStun(w http.ResponseWriter, r *http.Request) {
	register := r.URL.Query().Get("register")
	if register == "true" {
		uname := r.Context().Value("username")
		if uname == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(uname)
		peer, err := s.store.UpdateUserIp(r.Context().Value("username").(string), r.RemoteAddr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		data, err := json.Marshal(peer)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(data)

	} else {
		w.WriteHeader(200)
		w.Write([]byte(r.RemoteAddr))
	}
}
