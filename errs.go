package main

import (
	"net/http"
)

type ErrorPageS struct {
	ErrorCode    int
	ErrorMessage string
}

func panicIfErr(err error) {
	if err != nil {
		Log.PanicNF(1, "%s", err)
		panic(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	Log.DebugNF(1, "Redirecting to "+url)
	http.Redirect(w, r, url, 302)
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	SendErrCode(w, 404)
}

func SendErrCode(w http.ResponseWriter, code int) {
	pageDat := ErrorPageS{}
	pageDat.ErrorCode = code

	switch code {
	case 400:
		pageDat.ErrorMessage = "Má requisição."
	case 401:
		pageDat.ErrorMessage = "Não autorizado."
	case 404:
		pageDat.ErrorMessage = "Página não encontrada."
	case 403:
		pageDat.ErrorMessage = "Proibido. Você não tem a permissão necessária."
	case 409:
		pageDat.ErrorMessage = "Conflito."
	case 500:
		pageDat.ErrorMessage = "Erro interno do servidor."
	default:
		Log.Warning("Unknown http error code:", code)
		pageDat.ErrorMessage = "Erro desconhecido"
	}

	w.WriteHeader(code)

	w.Header().Set("Content-Type", "text/html")
	Templ.ExecuteTemplate(w, "error.html", pageDat)
}

func SendErrCodeAndLog(w http.ResponseWriter, code int, err interface{}) {
	Log.WarningNF(1, "Sending %d HTTP Error Code due to: %v", code, err)
	SendErrCode(w, code)
}
