package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pei0804/go-test-example/ui/session"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", index)
	r.Post("/login", login)
	r.Get("/member", member)
	http.ListenAndServe(":8080", r)
}

type Index struct {
	Msg string
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("name") == "" || r.FormValue("pw") != "pw" {
		errHandle(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	sess, _ := session.Get(r, "sess")
	sess.Values["name"] = r.FormValue("name")
	if err := session.Save(r, w, sess); err != nil {
		errHandle(w, r, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/member", http.StatusMovedPermanently)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		errHandle(w, r, http.StatusInternalServerError, err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		errHandle(w, r, http.StatusInternalServerError, err)
		return
	}
}

type Member struct {
	Name string
}

func member(w http.ResponseWriter, r *http.Request) {
	name := currentName(r)
	if name == "" {
		errHandle(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	tmpl, err := template.ParseFiles("member.html")
	if err != nil {
		errHandle(w, r, http.StatusInternalServerError, err)
		return
	}
	b := Member{Name: name}
	err = tmpl.Execute(w, b)
	if err != nil {
		errHandle(w, r, http.StatusInternalServerError, err)
		return
	}
}

func errHandle(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	io.WriteString(w, err.Error())
}

// CurrentName returns current user name who logged in.
func currentName(r *http.Request) string {
	if r == nil {
		return ""
	}
	sess, _ := session.Get(r, "sess")
	rawname, ok := sess.Values["name"]
	if !ok {
		return ""
	}
	name, ok := rawname.(string)
	if !ok {
		return ""
	}
	return name
}
