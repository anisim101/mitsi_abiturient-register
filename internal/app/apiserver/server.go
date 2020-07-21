package apiserver

import (
	"archive/zip"
	. "encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"io"
	"io/ioutil"
	"mitsoСhat/internal/app/model"
	"mitsoСhat/internal/app/store"
	"net/http"
	"os"
	"path/filepath"
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

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip  + file.Name() + "/")
		}
	}
}

func (s *server) Configure() {
	c := cron.New()
	c.AddFunc("0 0 * * *", RunEveryDay)
	c.Start()

}

func RemoveContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}
func RunEveryDay() {


	folders := [3]string{"gomel", "minsk", "vitebsk"}

	for _, folder := range folders {
		sourceFolder := "./files/" +  folder + "/"

		var time = time.Now().Format(time.RFC822)
		var url = "./arhives/" + folder + "_" + time + ".zip"
		// Get a Buffer to Write To
		outFile, err := os.Create(url)
		if err != nil {
			fmt.Println(err)
		}
		defer outFile.Close()

		// Create a new zip archive.
		w := zip.NewWriter(outFile)

		// Add some files to the archive.
		addFiles(w, sourceFolder, "")

		if err != nil {
			fmt.Println(err)
		}

		// Make sure to check the error on Close.
		err = w.Close()
		if err != nil {
			fmt.Println(err)
		}


		RemoveContents(sourceFolder)
	}
}

func (s *server) handleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.logRequest(w, r)

		uni := model.Abiturient{}

		if err := NewDecoder(r.Body).Decode(&uni); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		var baseUrl = "/home/uroot/abit_files/"

		if uni.UniverLocation == "Минск" {
			baseUrl = baseUrl + "minsk/"
		} else if uni.UniverLocation == "Витебск" {
			baseUrl = baseUrl + "vitebsk/"
		} else if uni.UniverLocation == "Гомель" {
			baseUrl = baseUrl + "gomel/"
		}

		baseUrl = baseUrl + uni.SerialAndPassportNumber + ".xml"

		if Exists(baseUrl) {
			s.error(w, r, http.StatusGone, store.FileExist)
			return
		}

		writer, er := os.Create(baseUrl)
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

		writer.Close()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// if err := s.store.AddPersson(&uni); err != nil {
		// 	s.error(w, r, 500, err)
		// 	return
		// }
	}
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
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
	s.router.PathPrefix("/abiturient_files/").Handler(http.StripPrefix("/abiturient_files/",
		http.FileServer(http.Dir(/*"./abiturient_files/" */"/home/uroot/abit_files/"))))
	s.router.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir(/*"./web/"*/ "/home/uroot/web/"))))

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
