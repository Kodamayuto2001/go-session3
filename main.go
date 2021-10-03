package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"fmt"
	"time"
	"html/template"
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"
)

var session_name string = "gsid"

var store *sessions.CookieStore

var session *sessions.Session 

type Data1 struct {
	Count	int
	Msg		string
}

func main() {
	gob.Register(&Data1{})

	sessionInit()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, session_name)

		data1, ok := session.Values["data1"].(*Data1)
		if data1 != nil {
			data1.Count++
			data1.Msg = fmt.Sprintf("%d件カウント", data1.Count)
		} else {
			data1 = &Data1{0, "データなし"}
		}
		fmt.Println(ok)
		fmt.Println(data1)
		session.Values["data1"] = data1
		
		sessions.Save(r, w)

		tmpl := template.Must(template.New("index").ParseFiles("templates/index.html"))
		tmpl.Execute(w, struct {
			Detail *Data1
		} {
			Detail: data1,
		})

		fmt.Print(time.Now())
		fmt.Println(" url = " + r.URL.Path)
	})

	r.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		sessionInit()
		
		fmt.Print(time.Now())
		fmt.Println(" url = " + r.URL.Path)

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.Handle("/", r)

	fmt.Println("localhost:3000")
	http.ListenAndServe(":3000", nil)
}

func sessionInit() {
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	str := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")

	store = sessions.NewCookieStore([]byte(str))
	session = sessions.NewSession(store, session_name)

	store.Options = &sessions.Options{
		Domain:		"localhost",
		Path:		"/",
		MaxAge:		0,
		Secure:		false,
		HttpOnly:	true,
	}

	fmt.Println("key		data --")
	fmt.Println(str)
	fmt.Println("")
	fmt.Println("store		data --")
	fmt.Println(store)
	fmt.Println("")
	fmt.Println("session	data --")
	fmt.Println(session)
	fmt.Println("")
}