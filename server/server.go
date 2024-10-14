package main

import (
	"database/sql"
	"log"
	// "mailinglist/jsonapi"
	"mailinglist/mdb"
	"sync"
    "mailinglist/grpcapi"
	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
	BindGrpc string `arg:"env:MAILINGLIST_BIND_GRPC"`
}

func main() {

	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8081"
	}
	if args.BindGrpc ==""{
		args.BindGrpc =":8080"
	}
	log.Printf("using database: '%v'\n", args.DbPath)

	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mdb.TryCreate(db)

	var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	log.Printf("starting JsonApi server\n")
	// 	jsonapi.Server(db, args.BindJson)
	// 	wg.Done()
	// }()
	// wg.Wait()



	wg.Add(1)
	go func() {
		log.Printf("starting Grpc server\n")
		grpcapi.Server(db, args.BindGrpc)
		wg.Done()
	}()
	wg.Wait()
}
