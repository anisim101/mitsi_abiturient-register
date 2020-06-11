package apiserver

import (
	. "encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"mitsoСhat/internal/app/model"
	"mitsoСhat/internal/app/store"
	"mitsoСhat/internal/app/store/sqlstore"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  sqlstore.Store
}

var (
	errorNotAuthorize           = errors.New("no authorization token")
	errIncorrectEmailOrPassword = errors.New("email or password not valid")
	TOKEN_SECRET_RULE           = "mitsoChatSecret"
)


func (s *server) handleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.logRequest(w,r)

		uni := model.Uni{}

		if err := NewDecoder(r.Body).Decode(&uni); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.AddPersson(&uni); err != nil {
			s.error(w,r, 500, err)
			return
		}
	}
}

func newServer(store sqlstore.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/api/upload_photo", s.handleUploadPhoto()).Methods("POST")
	s.router.HandleFunc("/api/addUser", s.handleAddUser()).Methods("POST")
	s.router.Handle("/files/photos/{rest}", http.StripPrefix("/files/photos/", http.FileServer(http.Dir("files/photos/"))))
}

func (s *server) logRequest(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	logger.Infof("started %s %s %s", r.Method, r.RequestURI, r.Header)
}



func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{
		"error": err.Error(),
	})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = NewEncoder(w).Encode(data)
	}
}

func (s *server) fetchToken(w http.ResponseWriter, r *http.Request) (string, error) {
	requestToken := r.Header.Get("Authorization")
	if requestToken == "" {
		return "", errorNotAuthorize
	}
	return requestToken, nil
}

func (s *server) handleUploadPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			ImageKey string `json:"image_key"`
		}
		s.logRequest(w, r)

		//10 mb limit
		r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
		file, handler, err := r.FormFile("file_name")
		if err != nil {
			s.error(w, r, http.StatusBadRequest, store.ErrBigFile)
			return
		}
		defer file.Close()

		fileName := handler.Filename
		split := strings.Split(fileName, ".")
		size := len(split)
		fileExisten := split[size-1]
		defer file.Close()

		savedFilePath := ""

		t := time.Now().Unix()
		savedFilePath = "files/photos/" + strconv.FormatInt(t, 10) + uuid.New().String() + "." + fileExisten
		f, err := os.OpenFile(savedFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, store.InternalError)
			return
		}
		defer f.Close()
		_, _ = io.Copy(f, file)
		s.respond(w, r, http.StatusOK, &response{
			ImageKey: savedFilePath,
		})
	}
}