package arbeit

import (
	"database/sql"
	"drive/domain"
	"time"
)

// Arbeit beschreibt eine Entitaet, fuer die Arbeitszeit erfasst werden soll.
type Job struct {
	name         string
	arbeitgeber  string
	arbeitnehmer string
	von          time.Time
	bis          time.Time
}

type Kalendertag struct {
	id       int
	jahr     int16
	monat    uint8
	tag      uint8
	datum    time.Time
	feiertag string
	//kw
	kw_jahr int
	kw_nr   uint8
	kw_tag  uint8

	yearday    uint16
	ordinal    int
	monatsname string
	tagesname  string
}

func (t Kalendertag) String() string {
	return ""
}
func (t Kalendertag) Gestern() {

}
func (t Kalendertag) Morgen() {

}

type Arbeitsjahr struct {
	ID                    int
	Account               *domain.Account
	Job                   *Job
	UrlaubVorjahr         sql.NullFloat64
	UrlaubAnspruch        sql.NullFloat64
	UrlaubTage            sql.NullFloat64
	UrlaubGeplant         sql.NullFloat64
	UrlaubRest            sql.NullFloat64
	Soll                  sql.NullFloat64
	Ist                   sql.NullFloat64
	Differenz             sql.NullFloat64
	tageFreizeitausgleich sql.NullFloat64
	tageKrank             sql.NullFloat64
	tageArbeit            sql.NullFloat64
	tageBuero             sql.NullFloat64
	tageDienstreise       sql.NullFloat64
	tageHomeoffice        sql.NullFloat64
	tageFrei              sql.NullFloat64
	jahrID                sql.NullInt64
	userID                sql.NullInt64
}

type Zeitspanne struct {
	Nr                                int
	Typ                               string
	Von, Bis                          *time.Time
	Dauer                             *float64
	Titel, Story, Beschreibung, Grund *string
	Arbeitszeit                       bool
}

type Arbeitstag struct {
	ID          int
	Account     domain.Account
	Job         Job
	Tag         Kalendertag
	Status      *string
	Typ         *string
	Soll        *float64
	Start       *time.Time
	Ende        *time.Time
	Brutto      *float64
	Pausen      *float64
	Netto       *float64
	Differenz   *float64
	Saldo       *float64
	Zeitspannen []Zeitspanne
}
