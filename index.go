package main

import (
    "bytes"
    "encoding/base64"
    "image"
    "image/jpeg"
    "io"
    "log"
    "net/http"
    "os"
    "text/template"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/show.html"))
/* template.ParseFiles: 別ファイルに定義したテンプレートを読み込む */
/* template.Must: バリデーションチェック */

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "index"}
    renderTemplate(w, "index", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    if err := templates.ExecuteTemplate(w, tmpl+".html", data); err != nil {
        log.Fatalln("Unable to execute template")
    }

}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // POSTのみをうける
    if r.Method != "POST" {
        // POST以外が来た場合はエラーを出す
        http.Error(w, "Allowed POST method only",
            http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(32 << 20) // maxMemory
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    file, _, err := r.FormFile("upload")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer file.Close()

    f, err := os.Create("/home/yoneda/github/go_lang/static/test.jpg")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer f.Close()

    io.Copy(f, file)
    http.Redirect(w, r, "/show", http.StatusFound)
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open("/home/yoneda/github/go_lang/static/test.jpg")
    // os.Open: 画像とりだす
    defer file.Close()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    img, _, err := image.Decode(file)
    // とりだしたファイルをデコードする
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    writeImageWithTemplate(w, "show", &img)
}

func writeImageWithTemplate(w http.ResponseWriter, tmpl string, img *image.Image) {
    buffer := new(bytes.Buffer)
    // 処理を行うためのメモリ領域を確保する
    if err := jpeg.Encode(buffer, *img, nil); err != nil {
        log.Fatalln("Unable to encode image.")
    }
    // エンコード処理
    str := base64.StdEncoding.EncodeToString(buffer.Bytes())
    data := map[string]interface{}{"Title": tmpl, "Image": str}
    renderTemplate(w, tmpl, data)
}

func main() {
    http.HandleFunc("/", IndexHandler)
    http.HandleFunc("/upload", UploadHandler)
    http.HandleFunc("/show", ShowHandler)

    err := http.ListenAndServe("192.168.0.69:8081", nil)
    if err != nil {
        log.Panic(err)
    }
}
