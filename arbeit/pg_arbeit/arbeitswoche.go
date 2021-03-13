package pg_arbeit

import (
	"drive/arbeit"
	"drive/errors"
	"time"

	_ "github.com/lib/pq"
)

func (d database) RetrieveArbeitstage(year, month, week int, accountID int) (a []arbeit.Arbeitstag, err error) {

	aa := []arbeit.Arbeitstag{}

	query := `
		SELECT id, status, typ, soll, start, ende, brutto, pausen, netto, differenz 
		  FROM go_arbeitstag a, kalendertag k  
		 WHERE a.tag_id=k.id
	`
	if month != 0 {
		query += "AND k.jahr_id=$1 AND k.monat=$2"
		err = d.DB.Select(&aa, query, year, month)
	}
	if week != 0 {
		query += "AND k.kw_jahr=$1 AND k.kw_jahr=$2"
		err = d.DB.Select(&aa, query, year, week)
	}
	if err != nil {
		err = errors.Wrap(err, "Could not Select  arbeitstage %v", aa)
	}

	return
}

type Arbeitstag2 struct {
	ID int
	//Account     domain.Account
	//Job         Job
	//Tag         Kalendertag
	Typ       *string // Arbeitstag, Wochenende, Feiertag
	Status    *string // Büro => 8, Dienstreise=>8, Krankheit=>0, Urlaub=>0, Zeitausgleich=>8
	Soll      *float64
	Start     *time.Time
	Ende      *time.Time
	Brutto    *float64
	Pausen    *float64
	Extra     *float64
	Netto     *float64 // Brutto + Extra - Pausen
	Differenz *float64 // Netto - Soll
	Saldo     *float64 // ergibt sich aus Saldo Vortag + Differenz
	//Zeitspannen []Zeitspanne
}
