package pg_arbeit

import (
	"database/sql"
	"drive/arbeit"
	"drive/errors"

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
			return nil, nil
		} else {
			return nil, errors.Wrap(err, "Could not QueryRow and Scan for Arbeitstag %v", id)
		}
	}

	pausen_query := `
		SELECT nr, typ, von, bis, dauer, titel, story, beschreibung, grund, arbeitszeit
		  FROM arbeits_zeitspanne 
		 WHERE arbeitstag_id = $1
	`
	zz := []arbeit.Zeitspanne{}
	err = d.DB.Select(&zz, pausen_query, id)
	if err != nil {
		return nil, errors.Wrap(err, "Could not Select  arbeits_zeitspanne %v", zz)
	}
	a.Zeitspannen = zz

	// rows, err := d.DB.Queryx(pausen_query, id)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Could not Query for arbeits_zeitspanne %v", id)
	// }
	// for rows.Next() {
	// 	var z arbeit.Zeitspanne
	// 	err = rows.StructScan(&z)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "Could not StructScan arbeits_zeitspanne")
	// 	}
	// 	a.Zeitspannen = append(a.Zeitspannen, z)
	// }

	return &a, nil
}
