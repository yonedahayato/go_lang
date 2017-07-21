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
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/show.html", "templates/upload.html", "templates/show_list.html"))
/* template.ParseFiles: 別ファイルに定義したテンプレートを読み込む */
/* template.Must: バリデーションチェック */

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "Index"}
    renderTemplate(w, "index", data)
}

func UploadDisplayHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "upload"}
    renderTemplate(w, "upload", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    if err := templates.ExecuteTemplate(w, tmpl+".html", data); err != nil {
        log.Fatalln("Unable to execute template. Template file is "+tmpl+".html")
    }
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // POSTのみをうける
    if r.Method != "POST" {
        // POST以外が来た場合はエラーを出す
        http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(32 << 20) // maxMemory
    CheckErr(err, w)

    file, handler, err := r.FormFile("upload")
    CheckErr(err, w)
    defer file.Close()

    DatabaseUpdate(handler.Filename, w)

    f, err := os.Create("paper/"+handler.Filename+".pdf")
    CheckErr(err, w)
    defer f.Close()

    io.Copy(f, file)
    http.Redirect(w, r, "/show_list", http.StatusFound)
}

func ShowListHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "show list"}
    renderTemplate(w, "show_list", data)
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

func DatabaseUpdate(filename string, w http.ResponseWriter) {
    if _, err := os.Stat("paper/paper.db"); err != nil{
        /*
        f, err := os.Create("paper/paper.db")
        CheckErr(err, w)
        defer f.Close()
        */

        db, err := sql.Open("sqlite3", "paper/paper.db")
        CheckErr(err, w)

        _, err = db.Exec("CREATE TABLE 'papers' ('id' INTEGER PRIMARY KEY AUTOINCREMENT, 'paper_name' TEXT)")
        CheckErr(err, w)

        fmt.Printf("i am here ========================================== \n")

        _, err = db.Exec("INSERT INTO 'papers' ('paper_name') VALUES (?)", filename)
        CheckErr(err, w)

        rows, err := db.Query("select * from 'papers'",)
        CheckErr(err, w)
        defer rows.Close()

        for rows.Next() {
            id, paper := 0, ""
            err = rows.Scan(&id, &paper)
            CheckErr(err, w)

            fmt.Println(id, paper)
        }
    } else {
        db, err := sql.Open("sqlite3", "paper/paper.db")
        CheckErr(err, w)

        _, err = db.Exec("INSERT INTO 'papers' ('paper_name') VALUES (?)", filename)
        CheckErr(err, w)

        rows, err := db.Query("select * from 'papers'",)
        CheckErr(err, w)
        defer rows.Close()

        for rows.Next() {
            id, paper := 0, ""
            err = rows.Scan(&id, &paper)
            CheckErr(err, w)

            fmt.Println(id, paper)
        }
    }
}

func CheckErr(err error, w http.ResponseWriter) {
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func main() {
    http.HandleFunc("/", IndexHandler)
    http.HandleFunc("/upload_display", UploadDisplayHandler)
    http.HandleFunc("/upload", UploadHandler)
    http.HandleFunc("/show", ShowHandler)
    http.HandleFunc("/show_list", ShowListHandler)

    err := http.ListenAndServe(":8082", nil)
    if err != nil {
        log.Panic(err)
        fmt.Errorf("err %s", err)
    }
}
