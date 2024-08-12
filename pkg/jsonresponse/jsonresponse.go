package jsonresponse

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, statuscode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(statuscode)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, errormsg string) {
	if code > 499 {
		log.Println(errormsg)
	}
	type errorResponse struct {
		Error string `json:"Error"`
	}

	RespondWithJson(w, code, errorResponse{Error: errormsg})

}
