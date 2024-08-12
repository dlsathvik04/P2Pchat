package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dlsathvik04/P2Pchat/db"
	"github.com/dlsathvik04/P2Pchat/pkg/hasher"
	"github.com/dlsathvik04/P2Pchat/pkg/jsonresponse"
	"github.com/dlsathvik04/P2Pchat/pkg/jwt"
)

type AuthAPIManager struct {
	hasher hasher.Hasher
	jwtman jwt.JWTManager
	db     *db.P2PchatDB
}

func NewAuthAPIManager(hasher hasher.Hasher, jwtman jwt.JWTManager, db *db.P2PchatDB) *AuthAPIManager {
	return &AuthAPIManager{hasher, jwtman, db}
}

func (ram *AuthAPIManager) Register(router *http.ServeMux) {
	router.HandleFunc("POST /register", ram.handleUserRegistration)
	router.HandleFunc("POST /login", ram.handleUserLogin)
}

func (ram *AuthAPIManager) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var registerBody struct {
		Username string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&registerBody)
	if err != nil {
		jsonresponse.RespondWithError(w, http.StatusBadRequest, "invalid body structure")
		return
	}
	peer, err := ram.db.CreateUser(registerBody.Username, ram.hasher.Hash(registerBody.Password), r.RemoteAddr)
	if err != nil {
		fmt.Println(err)
		jsonresponse.RespondWithError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	jsonresponse.RespondWithJson(w, http.StatusOK, peer)
}

func (ram *AuthAPIManager) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	// TODO - implement user login
	var loginBody struct {
		Username string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&loginBody)

	if err != nil {
		jsonresponse.RespondWithError(w, http.StatusBadRequest, "invalid body structure")
		return
	}

	peer, password, err := ram.db.GetUserByUsername(loginBody.Username)

	if err != nil {
		fmt.Println(err)
		jsonresponse.RespondWithError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	if ram.hasher.Compare(loginBody.Password, password) {
		token := ram.jwtman.GenerateToken(struct{ Username string }{peer.Username})
		jsonresponse.RespondWithJson(w, http.StatusOK, struct{ Token string }{token})
	} else {
		jsonresponse.RespondWithError(w, http.StatusBadRequest, "invalid credentials")
	}

}

func (ram *AuthAPIManager) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authtoken")
		if token == "" {
			jsonresponse.RespondWithError(w, http.StatusUnauthorized, "No token provided")
			return
		}
		var p struct{ Username string }
		err := ram.jwtman.AuthorizeToken(token, &p)

		if err != nil {
			jsonresponse.RespondWithError(w, http.StatusUnauthorized, "invalide identity")
			return
		}

		fmt.Println(p)
		ctx := context.WithValue(r.Context(), "username", p.Username)

		next(w, r.WithContext(ctx))
	}
}
