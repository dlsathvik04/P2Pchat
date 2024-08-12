package api

import (
	"net/http"

	"github.com/dlsathvik04/P2Pchat/db"
	"github.com/dlsathvik04/P2Pchat/pkg/jsonresponse"
)

type RendezvousAPIManager struct {
	store   db.P2PchatDB
	authman *AuthAPIManager
}

func NewRendezvousAPIManager(store db.P2PchatDB, aam *AuthAPIManager) *RendezvousAPIManager {
	return &RendezvousAPIManager{store, aam}
}

func (ram *RendezvousAPIManager) Register(server *http.ServeMux) {
	server.HandleFunc("GET /getip/", ram.authman.Authorize(ram.handleRendezvousRequest))
}

func (ram *RendezvousAPIManager) handleRendezvousRequest(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	peer, _, err := ram.store.GetUserByUsername(username)
	if err != nil {
		jsonresponse.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	jsonresponse.RespondWithJson(w, http.StatusOK, peer)

}
