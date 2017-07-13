package main

import (
    "bytes"
    "encoding/base64"
    "image"
    "image/jpeg"
    "io"
    "net/http"
    "os"
    "text/template"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))
/* template.ParseFiles: 別ファイルに定義したテンプレートを読み込む */
/* template.Must: バリデーションチェック */

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "index"}
    renderTemplate(w, "index", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    if err := templates.ExecuteTemplae(e, tmpl+".html", data); err != nil {
        log.Fatalln("Unable to execute template")
    }
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

}

func ShowHandler(w http.ResponseWriter, r *http.Request) {

}

func writeImageWithTemplate(w http.ResponseWriter, tmpl string, img *image.Image) {

}

func main() {
    http.HandleFunc("/", IndexHandler)
    http.ListenAndServe(":8081", nil)
}
