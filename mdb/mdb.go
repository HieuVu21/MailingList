package mdb

import (
	"database/sql"

	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type EmailEntry struct {
	Id        int64
	Email     string
	ConfirmAt *time.Time
	OptOut    bool
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
	 CREATE TABLE Emails (
	 id integer primary key ,
	 email text unique,
	 confirmed_at integer,
	 opt_out integer
	 );
	`)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code != 1 {
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}

	}
}
func EmailEntryFromRow(row *sql.Rows) (*EmailEntry, error) {
	var id int64
	var email string
	var ConfirmAt int64
	var OptOut bool

	err := row.Scan(&id, &email, &ConfirmAt, &OptOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t := time.Unix(ConfirmAt, 0)
	return &EmailEntry{Id: id, Email: email, ConfirmAt: &t, OptOut: OptOut}, nil
}
func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
INSERT INTO Emails (email, confirmed_at, opt_out)
VALUES (?, 0, false)`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetEmail(db *sql.DB, email string) (*EmailEntry, error) {
	rows, err := db.Query(`
	 SELECT id, email, confirmed_at, opt_out
        FROM Emails
        WHERE email = ?
	`, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return EmailEntryFromRow(rows)
	}
	return nil, nil
}

// func GetEmail(db *sql.DB, email string) (*EmailEntry, error) {
//     row := db.QueryRow(`
//      SELECT id, email, confirmed_at, opt_out
//         FROM Emails
//         WHERE email = ?
//     `, email)

//     var entry EmailEntry
//     err := row.Scan(&entry.Id, &entry.Email, &entry.ConfirmAt, &entry.OptOut)
//     if err != nil {
//         if err == sql.ErrNoRows {
//             return nil, fmt.Errorf("email not found")
//         }
//         log.Println(err)
//         return nil, err
//     }
//     return &entry, nil
// }

func UpdateEmail(db *sql.DB, entry *EmailEntry) error {
	t := entry.ConfirmAt.Unix()
	_, err := db.Exec(`
	 INSERT INTO Emails(email, confirmed_at, opt_out) 
VALUES (?, ?, ?) 
ON CONFLICT (email) DO UPDATE SET 
confirmed_at = ?, 
opt_out = ?`, entry.Email, t, entry.OptOut, t, entry.OptOut)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func DeleteEmai(db *sql.DB, email string) error {
	_, err := db.Exec(`
	update Emails 
	set opt_out = true 
	where email = ?
	`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type GetEmailBatchQueryParams struct {
	Page  int
	Count int
}

func GetEmailBatch(db *sql.DB, params GetEmailBatchQueryParams) ([]EmailEntry, error) {
	var empty []EmailEntry
	rows, err := db.Query(`
	select id,email,confirmed_at,opt_out from Emails 
	where opt_out =false 
	oder by id asc 
	limit ? offset ?
	
	`, params.Count, (params.Page-1)*params.Count)
	if err != nil {
		log.Println(err)
		return empty, err
	}
	defer rows.Close()
	emails := make([]EmailEntry, 0, params.Count)
	for rows.Next() {
		email, err := EmailEntryFromRow(rows)
		if err != nil {
			return nil, err
		}
		emails = append(emails, *email)
	}
	return emails, nil
}
