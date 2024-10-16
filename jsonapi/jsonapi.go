package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mailinglist/mdb"
	"net/http"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json; charset =utf-8")

}

func fromJson[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	json.Unmarshal(buf.Bytes(), &target)
}

	func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) {
		setJsonHeader(w)
		data, serverErr := withData()
		if serverErr != nil {
			w.WriteHeader(500)
			serverErrJson, err := json.Marshal(&serverErr)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(serverErrJson)
		}
		dataJson, err := json.Marshal(&data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Write(dataJson)
	}
// func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) {
// 	setJsonHeader(w)
// 	data, serverErr := withData()
// 	if serverErr != nil {

// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{"error": serverErr.Error()})
// 		return
// 	}

// 	json.NewEncoder(w).Encode(data)
// }

func returnErr(w http.ResponseWriter, err error, code int) {
	returnJson(w, func() (interface{}, error) {
		errMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(code)
		return errMessage, nil
	})
}

func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)
		if err := mdb.CreateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("Json create email %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}

func GetEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)

		returnJson(w, func() (interface{}, error) {
			log.Printf("Json getEmail %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

		})
	}
// func GetEmail(db *sql.DB) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		if req.Method != "GET" {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		email := req.URL.Query().Get("email")
// 		if email == "" {
// 			returnErr(w, errors.New("email parameter is required"), http.StatusBadRequest)
// 			return
// 		}

// 		returnJson(w, func() (interface{}, error) {
// 			log.Printf("Json getEmail %v\n", email)
// 			return mdb.GetEmail(db, email)
// 		})
// 	})
// }

func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)
		if err := mdb.UpdateEmail(db, &entry); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("Json UpdateEmail %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}
func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)
		if err := mdb.DeleteEmai(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("Json DeleteEmail %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}
func GetEmailBatch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			return
		}
		queryOptions := mdb.GetEmailBatchQueryParams{}
		fromJson(req.Body, &queryOptions)
		if queryOptions.Count <= 0 || queryOptions.Page <= 0 {
			returnErr(w, errors.New("error"), 400)
			return
		}
		returnJson(w, func() (interface{}, error) {
			log.Printf("Json GetEmailBatch: %v\n", queryOptions)
			return mdb.GetEmailBatch(db, queryOptions)
		})
	})
}

func Server(db *sql.DB, bind string) {
	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmail(db))
	http.Handle("/email/get_batch", GetEmailBatch(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))

	log.Printf("Starting server on %s", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatalf("Json server err: %v ", err)
	}

}
