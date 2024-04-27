package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ujjwal405/microservices/services/authentication/commonpkg/config"
	"github.com/stretchr/testify/require"
)

type MockQueue struct {
	hasErr error
}

func (mock MockQueue) AddToQueue(mail Mail) error {
	return mock.hasErr
}

type MockRepo struct {
	hasErr error
	isok   bool
	user   User
}

func (repo MockRepo) CountDocuments(email string) (bool, error) {
	return repo.isok, repo.hasErr
}
func (repo MockRepo) InsertUser(user User) error {
	return repo.hasErr
}
func (repo MockRepo) FetchUser(email string) (User, error) {
	return repo.user, repo.hasErr
}

type MockCache struct {
	isok bool
}

func (cache MockCache) AddClient(code string, detail *Details) {
	log.Println("adding client")
}
func (cache MockCache) Loop() {
	log.Println("looped")
}
func (cache MockCache) Check(code string) (Details, bool) {
	detail := Details{
		TTL: time.Now(),
		Detail: User{
			Email:    "abc125@gmail.com",
			Username: "abcdefgh",
			Password: "password123",
		},
	}
	return detail, cache.isok
}
func (cache MockCache) Awake() {
	log.Println("aweking")
}
func TestUserSignup(t *testing.T) {
	testcases := []struct {
		name          string
		body          User
		buildStub     func() *Handler
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "validationErr",
			body: User{
				Email:    "abcgmail.com",
				Username: "abcdef",
				Password: "cabd",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{}, MockCache{}, MockRepo{})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

			},
		},
		{
			name: "dbError",
			body: User{
				Email:    "abc125@gmail.com",
				Username: "abcdef",
				Password: "cabde",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{}, MockCache{}, MockRepo{isok: true, hasErr: errors.New("user already exists")})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
				data, err := io.ReadAll(recorder.Result().Body)
				require.NoError(t, err)
				defer recorder.Result().Body.Close()
				log.Println(string(data))

			},
		},
		{
			name: "QueueError",
			body: User{
				Email:    "abc125@gmail.com",
				Username: "abcdef",
				Password: "cabde",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{hasErr: errors.New("queue error")}, MockCache{}, MockRepo{})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
				data, err := io.ReadAll(recorder.Result().Body)
				require.NoError(t, err)
				defer recorder.Result().Body.Close()
				log.Println(string(data))

			},
		},

		{
			name: "NoError",
			body: User{
				Email:    "abc125@gmail.com",
				Username: "abcdef",
				Password: "cabde",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{}, MockCache{}, MockRepo{})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
				data, err := io.ReadAll(recorder.Result().Body)
				require.NoError(t, err)
				defer recorder.Result().Body.Close()
				log.Println(string(data))

			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			router := SetUpRouter()
			handler := tc.buildStub()
			router.POST("/signup", handler.Signup)
			res := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/signup", bytes.NewReader(data))
			require.NoError(t, err)
			router.ServeHTTP(res, req)
			tc.checkResponse(t, res)
		})
	}
}

func TestUserLogin(t *testing.T) {
	haspass := HashPassword("abcde")
	err := config.Loadconfig("../../..")
	require.NoError(t, err)
	testcases := []struct {
		name          string
		body          User
		buildStub     func() *Handler
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "DBError",
			body: User{
				Email:    "abc125@gmail.com",
				Username: "abcdef",
				Password: "cabde",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{}, MockCache{}, MockRepo{hasErr: errors.New("db error"), isok: false})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
				data, err := io.ReadAll(recorder.Result().Body)
				require.NoError(t, err)
				defer recorder.Result().Body.Close()
				log.Println(string(data))

			},
		},
		{

			name: "VerifyPassError",
			body: User{
				Email:    "abc125@gmail.com",
				Username: "abcdef",
				Password: "abcde",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{}, MockCache{}, MockRepo{isok: true, user: User{Email: "abc125@gmail.com", Username: "abcdef", Password: "cabde"}})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
				data, err := io.ReadAll(recorder.Result().Body)
				require.NoError(t, err)
				defer recorder.Result().Body.Close()
				log.Println(string(data))

			},
		},
		{
			name: "NoErr",
			body: User{
				Email:    "abc125@gmail.com",
				Username: "abcdef",
				Password: "abcde",
			},
			buildStub: func() *Handler {
				return MainHandler(MockQueue{}, MockCache{}, MockRepo{isok: true, user: User{Email: "abc125@gmail.com", Username: "abcdef", Password: haspass}})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
				data, err := io.ReadAll(recorder.Result().Body)
				require.NoError(t, err)
				defer recorder.Result().Body.Close()
				log.Println(string(data))
				header := recorder.Result().Header
				log.Println(header.Get("token"))
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			router := SetUpRouter()
			handler := tc.buildStub()
			router.POST("/login", handler.Login)
			res := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(data))
			require.NoError(t, err)
			router.ServeHTTP(res, req)
			tc.checkResponse(t, res)
		})
	}
}
