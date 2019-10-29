package utils

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"log"
	"net/http"
)

func SendError(err error, w http.ResponseWriter) {
	httpError, ok := err.(models.HTTPError)
	if !ok {
		w.WriteHeader(500) // return 500 Internal Server Error.
		return
	}

	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := httpError.ResponseHeaders() // GetUserByEmail http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)

	_, err = w.Write(body)

	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
}
