package arbeit

import "fmt"

type Repository interface {
	RetrieveArbeitstag(int, int, int, int) (*Arbeitstag, error)
	RetrieveArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error)
}

var Repo Repository

func GetArbeitstag(year, month, day int, accountID int) (*Arbeitstag, error) {

	a, err := Repo.RetrieveArbeitstag(year, month, day, accountID)
	fmt.Println("arbeitstag", a)
	return a, err
}

func GetArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error) {

	a, err := Repo.RetrieveArbeitsjahr(year, accountID)
	fmt.Println("ArbeitsJahr:", a)

	return a, err
}
