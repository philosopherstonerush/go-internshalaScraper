package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Dbase struct {
	db *sql.DB
}

func (d Dbase) add(of offer) {
	const INSERT_STATEMENT = `INSERT INTO offer (id, company, stipend, link) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`

	_, err := d.db.Exec(INSERT_STATEMENT, of.id, of.company, of.posted, of.stipend, of.link)
	if err != nil {
		log.Fatal(err)

	}
}

func (d Dbase) addAll(offers []offer) {
	INSERT_STATEMENT := "INSERT INTO offer (id, company, inserted_at, stipend, link) VALUES"
	const SUB = "('%s', '%s', '%s', '%s', '%s')"
	length := len(offers)
	for i := 0; i < length; i++ {
		str_value := fmt.Sprintf(SUB, offers[i].id, offers[i].company, offers[i].posted, offers[i].stipend, offers[i].link)
		INSERT_STATEMENT = INSERT_STATEMENT + " " + str_value
		if i+1 != length {
			INSERT_STATEMENT = INSERT_STATEMENT + ","
		}
	}
	const END_STATEMENT = " ON CONFLICT (id) DO NOTHING"

	INSERT_STATEMENT = INSERT_STATEMENT + END_STATEMENT
	_, err := d.db.Exec(INSERT_STATEMENT)
	if err != nil {
		log.Fatal(err)

	}
}
