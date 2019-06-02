package pg_arbeit

import (
	"database/sql"
	"drive/arbeit"
	"drive/errors"
	"log"

	_ "github.com/lib/pq"
)

func (d database) RetrieveArbeitstag(year, month, day int, accountID int) (*arbeit.Arbeitstag, error) {

	query := `
		SELECT id, status, typ, soll, start, ende, brutto, pausen, netto, differenz 
		  FROM arbeitstag 
		 WHERE id = $1
	`
	id := ((year*100+month)*100+day)*1000 + accountID

	a := arbeit.Arbeitstag{}
	err := d.DB.QueryRow(query, id).Scan(
		&a.ID,
		&a.Status,
		&a.Typ,
		&a.Soll,
		&a.Start,
		&a.Ende,
		&a.Brutto,
		&a.Pausen,
		&a.Netto,
		&a.Differenz,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// special case: there was no row
		} else {
			return nil, errors.Wrap(err, "Could not QueryRow and Scan for Arbeitstag %v", id)
		}
	}
	return &a, nil
}

//20151226001
func (d database) Asdf() {
	var (
		id   int
		name string
	)
	rows, err := d.DB.Query("select id, status from arbeitstag where id = $1", 20151226001)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
