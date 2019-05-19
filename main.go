package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gjvnq/go-logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/raja/argon2pw"

	lis "github.com/gjvnq/LibreLIS/libLIS"
)

var Config ConfigS
var Log *logger.Logger
var Templ *template.Template
var DB *sql.DB
var SessionStore *sessions.CookieStore
var SessionDuration time.Duration

type BasicPageDat struct {
	CurrentUser lis.User
}

func main() {
	var err error

	// Load config
	Log, err = logger.New("main", 1, os.Stdout)
	panicIfErr(err)
	Config = LoadConfigFile()

	// Set Logger
	Log.Levels["INFO"] = Config.DevMode
	Log.Levels["DEBUG"] = Config.DevMode
	if Config.DevMode {
		Log.Notice("Using development mode")
		pw, _ := argon2pw.GenerateSaltedHash("root")
		Log.Debug("Password 'root': " + pw)
	} else {
		Log.Notice("Using production mode")
	}
	Log.Info("SessionDuration = " + SessionDuration.String())

	// Prepare sessions
	SessionStore = sessions.NewCookieStore([]byte(Config.SessionSecret))

	// Connect to MySQL
	DB, err = sql.Open("mysql", Config.MySQL)
	lis.DB = DB
	lis.TheLogger = Log
	panicIfErr(err)
	defer DB.Close()

	// Load templates
	Templ, err = template.ParseGlob("./static/html/*.html")
	panicIfErr(err)

	// Prepare static file server
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))

	// Prepare router
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", LoginPage).Methods("GET", "POST")
	router.PathPrefix("/static/").Handler(static)
	router.NotFoundHandler = http.HandlerFunc(NotFoundPage)
	router.Use(RecoverPanics)
	router.Use(CheckLoggedIn)

	// Start server
	Log.Notice("Now listening... on " + Config.ListenOn + ":" + Config.Port)
	Log.Fatal(http.ListenAndServe(Config.ListenOn+":"+Config.Port, router))
}

func sendResponse(w http.ResponseWriter, mime string, data []byte) {
	var err error
	if mime != "" {
		w.Header().Set("Content-Type", mime)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		Log.ErrorNF(1, "Failed to write the response body: %v", err)
		return
	}
}

func sendTemplateResponse(w http.ResponseWriter, template_name string, template_data interface{}) {
	var err error

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	Templ.ExecuteTemplate(w, template_name, template_data)
	if err != nil {
		Log.ErrorNF(1, "Failed to write the response body: %v", err)
		return
	}
}

func GetString(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}

func RecoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(w http.ResponseWriter) {
			if r := recover(); r != nil {
				Log.StackAsError(fmt.Sprintf("Sending %d HTTP Error Code due to: %v", 500, r))
				SendErrCode(w, 500)
			}
		}(w)
		// Continue
		next.ServeHTTP(w, r)
	})
}
