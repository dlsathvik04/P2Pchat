package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dlsathvik04/P2Pchat/api"
	"github.com/dlsathvik04/P2Pchat/db"
	"github.com/dlsathvik04/P2Pchat/pkg/dotenv"
	"github.com/dlsathvik04/P2Pchat/pkg/hasher"
	"github.com/dlsathvik04/P2Pchat/pkg/jwt"
	_ "github.com/lib/pq"
)

func main() {

	dotenv.LoadDotEnv(".env", true)

	dburl := os.Getenv("DB_URL")
	serverSecret := os.Getenv("SECRET")
	hostName := os.Getenv("HOSTNAME")
	providerName := os.Getenv("PROVIDER_NAME")

	pgdb, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatal(err)
	}

	p2pchatDB := db.NewP2PchatDB(pgdb)
	// p2pchatDB.CleanUp()

	hasher := hasher.NewHasher(serverSecret)

	mux := http.NewServeMux()

	jwtman := jwt.NewJWTManager(time.Minute*10, serverSecret, providerName)

	server := api.NewAPIServer(hostName, mux, *p2pchatDB, hasher, jwtman)

	server.ListenAndServe()

	// stunManager

	// server := http.NewServeMux()
	// server.Handle("/stun/", stunserver)

	// http.ListenAndServe(":8000", server)
}
