package arbeit

import (
	"drive/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// A view function, or view for short (action), is simply a function that takes a Web request and returns a Web response.
// This response can be the HTML contents of a Web page, or a redirect, or a 404 error, or an XML document, or an image . . . or anything, really. The view itself contains whatever arbitrary logic is necessary to return that response. This code can live anywhere you want, as long as it’s on your Python path. There’s no other requirement–no “magic,” so to speak. For the sake of putting the code somewhere, the convention is to put views in a file called views.py, placed in your project or application directory.

func arbeit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)

	respond(w, r, "arbeit", map[string]interface{}{})
}

func arbeitsjahr(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	year, err := strconv.Atoi(ps.ByName("year"))
	if err != nil {
		errors.Error(w, r, err)
	}
	//sessionUser, _ := session.GetSessionUser(r, w)
	arbeitsjahr, err := GetArbeitsjahr(year, 1)
	if err != nil {
		errors.Error(w, r, err)
	}
	//fmt.Fprintf(w, "arbeitstag, %s!\n")
	respond(w, r, "arbeit", map[string]interface{}{
		"arbeitsjahr": arbeitsjahr,
	})
}

func arbeitsmonatwoche(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	year, err := strconv.Atoi(ps.ByName("year"))
	month, err2 := strconv.Atoi(ps.ByName("month"))
	if err != nil || err2 != nil {
		errors.Error(w, r, err)
		errors.Error(w, r, err2)
		return
	}
	//sessionUser, _ := session.GetSessionUser(r, w)
	arbeitstage, err := GetArbeitsMonat(year, month, 1)
	if err != nil {
		errors.Error(w, r, err)
		return
	}
	//fmt.Fprintf(w, "arbeitstag, %s!\n")
	respond(w, r, "arbeit", map[string]interface{}{
		"arbeitstage": arbeitstage,
	})
}

func arbeitstag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	year, err := strconv.Atoi(ps.ByName("year"))
	month, err := strconv.Atoi(ps.ByName("month"))
	day, err := strconv.Atoi(ps.ByName("day"))
	if err != nil {
		errors.Error(w, r, err)
		return
	}
	//sessionUser, _ := session.GetSessionUser(r, w)
	arbeitstag, err := GetArbeitstag(year, month, day, 1)
	if err != nil {
		errors.Error(w, r, err)
		return
	}
	//fmt.Fprintf(w, "arbeitstag, %s!\n")
	respond(w, r, "arbeit", map[string]interface{}{
		"arbeitstag": arbeitstag,
	})
}

func updateArbeitstag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	//sessionUser, _ := session.GetSessionUser(r, w)
	year, err := strconv.Atoi(ps.ByName("year"))
	month, err := strconv.Atoi(ps.ByName("month"))
	day, err := strconv.Atoi(ps.ByName("day"))
	if err != nil {
		errors.Error(w, r, err)
	}
	decoder := json.NewDecoder(r.Body)

	var a Arbeitstag
	err = decoder.Decode(&a)
	fmt.Println("a", a.Start, a.Ende)
	if err != nil {
		errors.Error(w, r, err)
	}
	a2, err := UpdateArbeitstag(year, month, day, 1, &a)

	fmt.Println("arbeitstag", a)

	//fmt.Fprintf(w, "arbeitstag, %s!\n")
	respond(w, r, "arbeit", map[string]interface{}{
		"arbeitstag": a2,
	})
}
