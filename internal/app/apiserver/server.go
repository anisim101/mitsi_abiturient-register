package apiserver

import (
	. "encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mitsoСhat/internal/app/model"
	"mitsoСhat/internal/app/store"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	//store  sqlstore.Store
}

var (
	errorNotAuthorize           = errors.New("no authorization token")
	errIncorrectEmailOrPassword = errors.New("email or password not valid")
	TOKEN_SECRET_RULE           = "mitsoChatSecret"
)

func (s *server) handleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.logRequest(w, r)

		uni := model.Abiturient{}

		if err := NewDecoder(r.Body).Decode(&uni); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		writer, er := os.Create("/home/uroot/abit_files/minsk/" + uni.SerialAndPassportNumber + ".xml")
		if er != nil {
			s.error(w, r, 500, er)
			return
		}

		encoder := xml.NewEncoder(writer)
		err := encoder.Encode(uni)

		if err != nil {
			s.error(w, r, 500, err)
			return
		}

		writer.Close();

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// if err := s.store.AddPersson(&uni); err != nil {
		// 	s.error(w, r, 500, err)
		// 	return
		// }
	}
}

func newServer( /*store sqlstore.Store*/ ) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		//store:  store,
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/api/get_user", s.handleGetUser()).Methods("POST")
	s.router.HandleFunc("/api/addUser", s.handleAddUser()).Methods("POST")
	s.router.Handle("/files/photos/{rest}", http.StripPrefix("/files/photos/", http.FileServer(http.Dir("./files/photos/"))))
	s.router.PathPrefix("/abiturient/").Handler(http.StripPrefix("/abiturient/",
		http.FileServer(http.Dir("./abiturient/"))))
		s.router.PathPrefix("/abiturient_files/").Handler(http.StripPrefix("/abiturient_files/",
		http.FileServer(http.Dir("/home/uroot/abit_files/"))))
}

func (s *server) logRequest(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	logger.Infof("started %s %s %s", r.Method, r.RequestURI, r.Header)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{
		"error": err.Error(),
	}, true)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, json bool) {
	w.WriteHeader(code)
	if data != nil {

		if json {
			_ = NewEncoder(w).Encode(data)
		} else {
			_ = xml.NewEncoder(w).Encode(data)
		}

	}
}

func (s *server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqest struct {
			PasportID string `json:"pasport_id"`
			IsJson    bool   `json:"is_json"`
		}
		s.logRequest(w, r)

		req := reqest{}

		if err := NewDecoder(r.Body).Decode(&req); err != nil {
			return
		}

		if req.PasportID == "" {
			//s.error(w, r, http.StatusBadRequest, Error)
			return
		}

		//err, person := s.store.GetPerson(req.PasportID)

		// if err != nil {
		// 	s.error(w, r, 500, err)
		// 	return
		// }

		w.Header().Set("Access-Control-Allow-Origin", "*")
		//s.respond(w, r, 200, person, req.IsJson)

	}
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.respond(w, r, http.StatusOK, &response{
			ImageKey: savedFilePath,
		}, true)
	}
}
