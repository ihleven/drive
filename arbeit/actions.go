package arbeit

import (
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
		fmt.Println(err)
	}
	//sessionUser, _ := session.GetSessionUser(r, w)
	arbeitsjahr, err := GetArbeitsjahr(year, 1)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Fprintf(w, "arbeitstag, %s!\n")
	respond(w, r, "arbeit", map[string]interface{}{
		"arbeitsjahr": arbeitsjahr,
	})
}

func arbeitstag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	year, err := strconv.Atoi(ps.ByName("year"))
	month, err := strconv.Atoi(ps.ByName("month"))
	day, err := strconv.Atoi(ps.ByName("day"))
	if err != nil {
		fmt.Println(err)
	}
	//sessionUser, _ := session.GetSessionUser(r, w)
	arbeitstag, err := GetArbeitstag(year, month, day, 1)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Fprintf(w, "arbeitstag, %s!\n")
	respond(w, r, "arbeit", map[string]interface{}{
		"arbeitstag": arbeitstag,
	})
}
