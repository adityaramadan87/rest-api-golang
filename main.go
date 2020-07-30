package main

import (
	"belajar-golang/routes/socket"
	"belajar-golang/routes/web"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Halo !! \n  Apa Kabar!!")
}

func main() {

	//set routes for HTTP
	http.Handle("/", web.SetRoutes())

	//for web
	http.HandleFunc("/index", index)

	//set Routes for SOCKET
	http.Handle("/socket.io/", socket.SetRoutes())
	fmt.Println("backend start server on port :9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}

	// log.Fatal(http.ListenAndServe(":9000", nil))

}
