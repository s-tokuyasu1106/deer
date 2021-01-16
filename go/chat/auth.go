package main

import (
	"net/http"
	"strings"
	"fmt"
	"log"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らか別のエラー
		panic(err.Error())
	} else {
		// 成功　ラップされたハンドラをよびだします
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}


// loginHandlerはサードパーティーへのログインの処理を受け持ちます
// パスの形式；/auth/{action}/{privider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeingAuthURLの呼び出し中にエラーが発生しました", provider, "-", err)
		}
		w.Header().Set("location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証の取得に失敗しました", provider, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした", provider, "-", err)
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーの取得に失敗しました", provider, "-", err)
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/"})
			w.Header()["Location"] = []string{"/chat"}
			w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}