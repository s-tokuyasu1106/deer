package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"example.com/m/v2/go/trace"
)

// temp1は1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

// serveHttpはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈します

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// ルート
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: login.html})
	// http.Handle("/", &templateHandler{filename: "chat.html"})

	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()
	// webサーバを開始します
	log.Println("webサーバを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
