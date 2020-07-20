package main

import (
	"database/sql"
	"log"
	. "mitsoСhat/internal/app/apiserver"

)


//var (
//	configPath string
//)

//func init() {
//	flag.StringVar(&configPath, "config-path",
//		"./config/apiserver.toml", "path to config file")
//}

func main() {
	// flag.Parse()
	config := NewConfig()
	// _, err := toml.DecodeFile(configPath, config)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// http.Handle("/abiturient/", http.StripPrefix("/abiturient/", http.FileServer(http.Dir("abiturient"))))
	// http.ListenAndServe(":50", nil)
	

	// db,err := newDB(config.DatabaseURL)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//defer db.Close()

	//store := sqlstore.New(db)

	if err := Start(config/*, *store*/); err != nil {
		log.Fatal(err)
	}

}

func newDB(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres",databaseUrl)

	if err != nil {
		return  nil,err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return  db, nil
}

