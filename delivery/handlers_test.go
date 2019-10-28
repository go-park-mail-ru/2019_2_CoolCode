package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Body       interface{}
	SessionID  string
	Headers    map[string]string
	Method     string
	URL        string
	Response   string
	StatusCode int
}

type InvalidJson struct {
	data string
}

var api = UserHandlers{
	Users:    useCase.NewUserUseCase(repository.NewArrayUserStore()),
	Sessions: make(map[string]uint64, 0),
}

var globalSessionID string

func AddContext(r *http.Request, key string, value string) {
	context.Set(r, key, value)
}

//func TestPreTest(t *testing.T) {
//	reader, _ := os.Open("users.txt")
//	defer reader.Close()
//	var users Users
//	decoder := json.NewDecoder(reader)
//	_ = decoder.Decode(&users)
//	api.Users.readUsers(users)
//}

func TestSignUp(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Body: models.User{
				Email:    "test1@test.com",
				Password: "1",
				Username: "test",
			},
			Method:     "POST",
			URL:        "/users",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Body: models.User{
				Email:    "test1@test.com",
				Password: "1",
			},
			Method:     "POST",
			URL:        "/users",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Body: InvalidJson{
				data: "mem",
			},
			Method:     "POST",
			URL:        "/users",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusOK,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		api.SignUp(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}
	}

}

func TestLogin(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Body: models.User{
				Email:    "test1@test.com",
				Password: "1",
			},
			Method:     "POST",
			URL:        "/login",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Body: models.User{
				Email:    "test1@test.com",
				Password: "2",
			},
			Method:     "POST",
			URL:        "/login",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Body: models.User{
				Email:    "test2@test.com",
				Password: "2",
			},
			Method:     "POST",
			URL:        "/login",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusBadRequest,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		api.Login(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}

		if w.Code == http.StatusOK {
			globalSessionID = w.Header().Get("Set-Cookie")
		}
	}

}

func TestSession(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Method:     "GET",
			SessionID:  globalSessionID,
			URL:        "/users",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Method:     "GET",
			URL:        "/users",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusUnauthorized,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range test.Headers {
			req.Header.Set(k, v)
		}
		if test.SessionID != "" {
			req.Header.Set("Cookie", test.SessionID)
		}

		w := httptest.NewRecorder()
		api.GetUserBySession(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}
	}

}

func TestGetUser(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Method:     "GET",
			SessionID:  globalSessionID,
			URL:        "/users/1",
			Response:   `{"id":1,"username":"test","email":"test1@test.com","fullname":"John Doe","password":"","fstatus":"","phone":""}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Method:     "GET",
			URL:        "/users/1",
			StatusCode: http.StatusUnauthorized,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range test.Headers {
			req.Header.Set(k, v)
		}
		if test.SessionID != "" {
			req.Header.Set("Cookie", test.SessionID)
		}

		w := httptest.NewRecorder()
		api.GetUser(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}

		resp := w.Result()
		respBody, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(respBody)
		if bodyStr != test.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				testNum, bodyStr, test.Response)
		}
	}

}

func TestFindUsers(t *testing.T) {
	cases := []TestCase{
		TestCase{
			SessionID:  globalSessionID,
			Method:     "GET",
			URL:        "/users/tes",
			StatusCode: http.StatusOK,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range test.Headers {
			req.Header.Set(k, v)
		}
		if test.SessionID != "" {
			req.Header.Set("Cookie", test.SessionID)
		}

		w := httptest.NewRecorder()
		api.FindUsers(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}
	}

}

func TestEditUser(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Body: models.User{
				ID:    1,
				Email: "test1@test.com",
			},
			Method:     "PUT",
			SessionID:  globalSessionID,
			URL:        "/users/1",
			StatusCode: http.StatusOK,
		},
		TestCase{
			Body: models.User{
				ID:       1,
				Email:    "test1@test.com",
				Password: "1",
			},
			Method:     "PUT",
			URL:        "/users/1",
			StatusCode: http.StatusUnauthorized,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range test.Headers {
			req.Header.Set(k, v)
		}
		if test.SessionID != "" {
			req.Header.Set("Cookie", test.SessionID)
		}

		w := httptest.NewRecorder()
		api.EditProfile(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}
	}

}

func TestLogout(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Method: "POST",

			SessionID:  globalSessionID,
			URL:        "/logout",
			StatusCode: http.StatusOK,
		},
		TestCase{
			Method: "POST",

			URL:        "/logout",
			StatusCode: http.StatusUnauthorized,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range test.Headers {
			req.Header.Set(k, v)
		}
		if test.SessionID != "" {
			req.Header.Set("Cookie", test.SessionID)
		}

		w := httptest.NewRecorder()
		api.Logout(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}
	}

}

func TestEditUserAfterLogout(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Body: models.User{
				ID:       1,
				Email:    "test1@test.com",
				Password: "1",
			},
			Method:     "PUT",
			SessionID:  globalSessionID,
			URL:        "/users/1",
			StatusCode: http.StatusUnauthorized,
		},
	}

	for testNum, test := range cases {
		userJSON, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewReader(userJSON)
		req, err := http.NewRequest(test.Method, test.URL, body)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range test.Headers {
			req.Header.Set(k, v)
		}
		if test.SessionID != "" {
			req.Header.Set("Cookie", test.SessionID)
		}

		w := httptest.NewRecorder()
		api.EditProfile(w, req)

		if w.Code != test.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				testNum, w.Code, test.StatusCode)
		}
	}

}

//func TestGetPgoto(t *testing.T) {
//	cases := []TestCase{
//		TestCase{
//			Method:     "GET",
//			URL:        "/photos/1",
//			StatusCode: http.StatusUnauthorized,
//		},
//	}
//
//	for testNum, test := range cases {
//		userJSON, err := json.Marshal(test.Body)
//		if err != nil {
//			t.Fatal(err)
//		}
//		body := bytes.NewReader(userJSON)
//		req, err := http.NewRequest(test.Method, test.URL, body)
//		if err != nil {
//			t.Fatal(err)
//		}
//		for k, v := range test.Headers {
//			req.Header.Set(k, v)
//		}
//		if test.SessionID != "" {
//			req.Header.Set("Cookie", test.SessionID)
//		}
//
//		w := httptest.NewRecorder()
//		api.getPhoto(w, req)
//
//		if w.Code != test.StatusCode {
//			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
//				testNum, w.Code, test.StatusCode)
//		}
//	}
//
//}
