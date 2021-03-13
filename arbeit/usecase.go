package arbeit

import (
	"drive/errors"
	"fmt"
)

type Repository interface {
	ListArbeitstage(year, month, week int, accountID int) ([]Arbeitstag, error)
	RetrieveArbeitstag(int, int, int, int) (*Arbeitstag, error)
	RetrieveArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error)
	UpdateArbeitstag(int, *Arbeitstag) error
	ListZeitspannen(int) ([]Zeitspanne, error)
	UpsertZeitspanne(int, *Zeitspanne) error
	DeleteZeitspanne(*Zeitspanne) error
}

var Repo Repository

func GetArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error) {

	a, err := Repo.RetrieveArbeitsjahr(year, accountID)
	fmt.Println("ArbeitsJahr:", a)

	return a, err
}

func GetArbeitsMonat(year, month, accountID int) ([]Arbeitstag, error) {
	fmt.Println("arbeitsmonat", year, month)
	aa, err := Repo.ListArbeitstage(year, month, 0, accountID)

	return aa, err
}

func GetArbeitstag(year, month, day int, accountID int) (*Arbeitstag, error) {
	fmt.Println("arbeitstag")
	a, err := Repo.RetrieveArbeitstag(year, month, day, accountID)
	fmt.Println("arbeitstag", a)
	return a, err
}

func UpdateArbeitstag(year, month, day int, accountID int, arbeitstag *Arbeitstag) (*Arbeitstag, error) {

	a, err := Repo.RetrieveArbeitstag(year, month, day, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "Could not retrieve Arbeitstag %s%s%s", year, month, day)
	}
	if arbeitstag.Start != nil && arbeitstag.Ende != nil {
		arbeitstag.Brutto = arbeitstag.Ende.Sub(*arbeitstag.Start).Hours()
		arbeitstag.Netto = arbeitstag.Brutto - arbeitstag.Pausen + arbeitstag.Extra
		arbeitstag.Differenz = arbeitstag.Soll - arbeitstag.Netto
	}
	err = UpdateZeitspannen(a.ID, arbeitstag.Zeitspannen)

	Repo.UpdateArbeitstag(a.ID, arbeitstag)

	return arbeitstag, err
}

func UpdateZeitspannen(arbeitstagId int, zeitspannen []Zeitspanne) error {
	fmt.Println("UpdateZeitspanne")
	// list of current zeitspanne ids
	//zeitspanne_ids := make([]int, 0)
	///for _, zeitspanne := range zeitspannen {
	//	zeitspanne_ids = append(zeitspanne_ids, zeitspanne.Nr)
	//}
	dbZeitspannen, _ := Repo.ListZeitspannen(arbeitstagId)
	for _, zeitspanne := range dbZeitspannen {
		if !IsContained(zeitspannen, zeitspanne) {
			Repo.DeleteZeitspanne(&zeitspanne)
		}
	}
	for _, zeitspanne := range zeitspannen {
		err := Repo.UpsertZeitspanne(arbeitstagId, &zeitspanne)
		if err != nil {
			fmt.Println("asdfasdfasdf", err)
		}
	}

	//if not zeitspanne.get('nr', False):
	//    max_zeitspanne_num += 1
	//    zeitspanne['nr'] = max_zeitspanne_num

	//zeitspanne, created = Zeitspanne.objects.update_or_create(arbeitstag=instance, nr=zeitspanne['nr'], defaults=zeitspanne)
	//zeitspanne.eval()
	return nil
}

func IsContained(haystack []Zeitspanne, needle Zeitspanne) bool {
	for _, n := range haystack {
		if n.Nr == needle.Nr {
			return true
		}
	}
	return false
}
