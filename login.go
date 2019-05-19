package main

import (
	"net/http"
	"strings"
	"time"

	lis "github.com/gjvnq/LibreLIS/libLIS"
	"github.com/gorilla/context"
)

type Key string

const UserKey Key = "UserKey"

type LoginPageS struct {
	Err   string
	Email string
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, err := SessionStore.Get(r, "session")
	if err != nil {
		Log.Warning(err)
	}
	session.Values["user"] = nil
	err = session.Save(r, w)
	panicIfErr(err)
	redirect(w, r, "/")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	page_data := LoginPageS{}

	// Check logged in
	user := GetCurrentUser(r)
	if user != nil {
		redirect(w, r, "/home")
		return
	}

	if r.FormValue("email") != "" {
		user := lis.LoadUserByEmail(r.FormValue("email"))
		if user == nil {
			page_data.Err = "Usuário não encontrado: " + r.FormValue("email")
		} else {
			if !user.VerifyPassword(r.FormValue("password")) {
				page_data.Err = "Senha incorreta"
			} else {
				// Write session
				session, err := SessionStore.Get(r, "session")
				panicIfErr(err)
				session.Values["user_id"] = user.Id
				session.Values["valid_until"] = time.Now().Add(SessionDuration)
				err = session.Save(r, w)
				panicIfErr(err)
				redirect(w, r, "/home")
				return
			}
		}
	}
	page_data.Email = r.FormValue("email")
	sendTemplateResponse(w, "login.html", page_data)
}

func GetCurrentUser(r *http.Request) *lis.User {
	val, ok := context.Get(r, UserKey).(*lis.User)
	if ok {
		return val
	}
	return nil
}

func CheckLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// No need to be logged in for these pages:
		if r.RequestURI == "/login" || strings.HasPrefix(r.RequestURI, "/static/") {
			next.ServeHTTP(w, r)
			return
		}
		// Clear just in case
		context.Set(r, UserKey, nil)
		// Get session
		session, err := SessionStore.Get(r, "session")
		panicIfErr(err)
		// Check if time has expired
		valid_until, ok := session.Values["valid_until"].(time.Time)
		// Redirect non logged in users
		if !ok || time.Now().After(valid_until) {
			redirect(w, r, "/login")
			return
		} else {
			// Update session expiration
			session.Values["valid_until"] = time.Now().Add(SessionDuration)
			err = session.Save(r, w)
			panicIfErr(err)
		}
		// Load user id from session
		user_id, ok := session.Values["user_id"].(int)
		if !ok {
			redirect(w, r, "/login")
			return
		}
		// Load the actual user "object"
		user := lis.LoadUserById(user_id)
		if user == nil {
			redirect(w, r, "/login")
			return
		}
		// Save user "object" so we can access it later
		context.Set(r, UserKey, user)
		// User logged in, continue
		next.ServeHTTP(w, r)
	})
}
