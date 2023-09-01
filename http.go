package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Простейший веб-сервер
const portNumber = ":8000"

func UIDPage(w http.ResponseWriter, r *http.Request) {
	// Проверка на то, что длина url равна 19 - длине UID
	if r.URL.Path != "/" {
		if len(r.URL.Path[1:]) != 19 || r.Method != "GET" {
			http.Error(w, "404 not found.", http.StatusNotFound)
		} else {
			order, exist := GetOrder(r.URL.Path[1:])
			w.Header().Set("Content-Type", "application/json")
			if exist {
				var data []byte
				data, _ = json.Marshal(order)
				fmt.Fprint(w, string(data))
			} else {
				fmt.Fprint(w, "{}")
			}
		}
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "data/form.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		UID := r.FormValue("UID")
		http.Redirect(w, r, "/"+UID, http.StatusSeeOther)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func StartHttp() error {
	http.HandleFunc("/", UIDPage)
	if err := http.ListenAndServe(portNumber, nil); err != http.ErrServerClosed {
		return err
	}
	return nil
}
 
