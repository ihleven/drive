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
	//id    int
	Jahr     int16 `db:"jahr_id" json:"year"`
	Monat    uint8 `db:"monat" json:"month"`
	Tag      uint8 `db:"tag" json:"day"`
	Datum    time.Time
	Feiertag *string
	KwJahr   int   `db:"kw_jahr" json:"kw_jahr"`
	KwNr     uint8 `db:"kw_nr" json:"kw_nr"`
	KwTag    uint8 `db:"kw_tag" json:"kw_tag"`

	Jahrtag uint16 `db:"jahrtag" json:"jahrtag"`
	Ordinal int
	//monatsname string
	//tagesname  string
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
	Account               *domain.Account `json:"-"`
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
	ID           int `db:"id" json:"id"`
	Account      domain.Account
	Job          Job
	Status       string `db:"status" json:"status"`
	Kategorie    string `db:"kategorie" json:"kategorie"`
	Krankmeldung bool
	Urlaubstage  float64
	Soll         float64
	Start        *time.Time `db:"beginn" json:"beginn"`
	Ende         *time.Time `db:"ende" json:"ende"`
	Brutto       float64
	Pausen       float64
	Extra        float64
	Netto        float64
	Differenz    float64
	Saldo        *float64
	Zeitspannen  []Zeitspanne
	Kalendertag  ` json:"kalendertag"`
}
