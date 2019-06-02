package arbeit

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterSubRouter(register func(string, func(http.ResponseWriter, *http.Request))) {

	arbeitRoutingPrefix := "/arbeit/"

	router := httprouter.New()
	router.GET("/", arbeit)
	router.GET("/:year", arbeitsjahr)
	router.GET("/:year/:month/:day", arbeitstag)

	arbeitRoutingFunc := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = r.URL.Path[len(arbeitRoutingPrefix)-1:]
		fmt.Println("routing arbeitsroute => ", r.URL.Path)
		router.ServeHTTP(w, r)
	}
	register(arbeitRoutingPrefix, arbeitRoutingFunc)

}
