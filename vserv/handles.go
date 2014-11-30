package main

import (
	"bitbucket.org/strings/via/pkg"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/plan/", plan)
}

func root(w http.ResponseWriter, r *http.Request) {
	html(w)
	pf, err := via.PlanFiles()
	if err != nil {
		fmt.Fprintln(w, err)
	}
	plans := []*via.Plan{}
	for _, f := range pf {
		p, e := via.ReadPath(f)
		if e != nil {
			log.Fatal(err)
		}
		plans = append(plans, p)
	}
	t, e := template.ParseFiles("root.tmpl")
	if e != nil {
		fmt.Fprintln(w, e)
	}
	e = t.Execute(w, plans)
	if e != nil {
		fmt.Fprintln(w, e)
	}
}

func plan(w http.ResponseWriter, r *http.Request) {
	html(w)
	t, e := template.ParseFiles("plan.tmpl")
	if e != nil {
		fmt.Fprintln(w, e)
	}
	p, e := via.FindPlan(path.Base(r.URL.String()))
	if e != nil {
		fmt.Fprintln(w, e)
	}
	p.Flags = via.GetConfig().Flags
	e = t.Execute(w, p)
	if e != nil {
		fmt.Fprintln(w, e)
	}
}

func html(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}
