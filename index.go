package main

import (
    "log"
    "net/http"
    "time"
    "html/template"
)

func clockHandler(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles("/home/yoneda/github/go_lang/templates/clock.html.tpl"))
    if err := t.ExecuteTemplate(w, "clock.html.tpl", time.Now()); err != nil {
        log.Fatal(err)
    }
}

func main() {
    http.HandleFunc("/clock/", clockHandler)
    http.Handle("/static",
        http.StripPrefix("/static/",
            http.FileServer(http.Dir("/home/yoneda/github/go_lang"))))

    log.Fatal(
        http.ListenAndServe(":8081", nil))
}
