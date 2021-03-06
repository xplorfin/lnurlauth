package integration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	lnurlHelper "github.com/fiatjaf/go-lnurl"
	"github.com/xplorfin/lnurlauth"
	"github.com/xplorfin/lnurlauth/storage"
)

const CookieName = "lnurlAuth"

var sessionStore storage.MemorySessionStore

func ParseUrl(rawUrl string) lnurlHelper.LNURLAuthParams {
	parsed, _ := url.Parse(rawUrl)
	params, _ := lnurlHelper.HandleAuth(rawUrl, parsed, parsed.Query())
	return params.(lnurlHelper.LNURLAuthParams)
}

func isAuthenticated(w http.ResponseWriter, r *http.Request) (isAuthenticated bool) {
	authToken := storage.CookieStore(w, r).Get(CookieName)
	if authToken != "" {
		authParams := ParseUrl(authToken)
		if authParams.K1 != "" {
			sessionData := sessionStore.Get(authParams.K1)
			if sessionData != nil && sessionData.Key != "" {
				isAuthenticated = true
			}
		}
	}
	return isAuthenticated
}

func returnJson(v interface{}, w http.ResponseWriter) {
	res, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func GenerateServer() http.Server {
	res := http.NewServeMux()

	// redirect to login page
	res.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("local")
		ok := isAuthenticated(w, r)
		HomeTpl.Execute(w, ok)
	})

	// return true if a user is authenticated, false otherwise
	res.HandleFunc("/is-authenticated", func(w http.ResponseWriter, r *http.Request) {
		status := AuthStatus{
			IsAuthenticated: isAuthenticated(w, r),
		}
		returnJson(status, w)
	})

	res.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		storage.CookieStore(w, r).Remove(CookieName)
		http.Redirect(w, r, "/", 302)
	})


	// redirect to login page
	res.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if isAuthenticated(w, r) {
			http.Redirect(w, r, "/", 302)
			return
		}
		authToken := storage.CookieStore(w, r).Get(CookieName)
		var encodedUrl, parsedUrl string
		if authToken == "" {
			encodedUrl, parsedUrl, _ = lnurlauth.GenerateLnUrl(fmt.Sprintf("http://%s/%s", r.Host, "callback"))
			http.SetCookie(w, &http.Cookie{Name: CookieName, Value: parsedUrl, HttpOnly: false})
			sessionStore.Set(ParseUrl(parsedUrl).K1, lnurlauth.SessionData{
				LnUrl:  encodedUrl,
				Key:    "",
			})
		} else {
			encodedUrl, _ = lnurlHelper.LNURLEncode(authToken)
		}

		qrCode, _ := lnurlauth.GenerateQrCode(encodedUrl)
		qrString := base64.StdEncoding.EncodeToString(qrCode)

		LoginPage.Execute(w, LoginPageData{
			Encoded:   encodedUrl,
			DataUri:   template.URL(fmt.Sprintf("data:image/png;base64,%s", qrString)),
			CancelUrl: "",
		})
	})

	res.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		key, k1, err := lnurlauth.Authenticate(r)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		sessionData := sessionStore.Get(k1)
		if sessionData != nil {
			sessionData = &lnurlauth.SessionData{}
		}

		sessionStore.Set(k1, lnurlauth.SessionData{
			LnUrl:  sessionData.LnUrl,
			Key:    key,
		})

		returnJson(CallbackStatus{Ok: true}, w)
	})



	return http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res.ServeHTTP(w, r)
		}),
	}
}
