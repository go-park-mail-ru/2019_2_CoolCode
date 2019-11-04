package delivery

import (
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
)

var usersApi UserHandlers

func Initial(t *testing.T) {

}
func TestUserHandlers_SignUp(t *testing.T) {

	apitest.New().
		Handler(http.HandlerFunc(usersApi.SignUp)).
		Post("users").
		Body(`{`)
}
