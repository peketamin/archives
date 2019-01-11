package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func serve() {
	// after db
	flag.Parse()
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/image/", imageHandler)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", viewHandler)
	// handle static
	http.Handle("/"+staticRoot+"/", http.StripPrefix("/"+staticRoot, http.FileServer(http.Dir("./"+staticRoot+"/"))))

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}
	http.ListenAndServe(":8080", nil)
}

func main() {
	dbconnect("db.dat")
	_init()
	serve()
}
