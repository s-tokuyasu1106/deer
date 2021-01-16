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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.temp1.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈します

	// Comniauthのセットアップ
	gomniauth.SetSecurityKey("AIzaSyDCjcQz1LII1znGD7agarraagaararaghaNHAMluBjKv89xpvCs")
	gomniauth.WithProviders(
		google.New("684305737811-54nfm6jvgp1i7tnlf9fil5mj1sk872bp.apps.googleusercontent.com", "wXBdxRR8CZic0mL9jl9Aj3at", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// ルート
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
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
