package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"path"
	"text/template"
	"time"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	creds := model.Credentials{} // TODO: replace this
	username := r.FormValue("username")
	password := r.FormValue("password")

	creds.Username = username
	creds.Password = password

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
		return
	}
	// jika ok
	w.WriteHeader(200)

	err := api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{"name": creds.Username, "message": "register success!"}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	// Read usernmae and password request with FormValue.
	creds := model.Credentials{} // TODO: replace this
	// err := json.NewDecoder(r.Body).Decode(&creds)
	username := r.FormValue("username")
	password := r.FormValue("password")

	creds.Username = username
	creds.Password = password

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
		return
	}

	// TODO: answer here

	err := api.usersRepo.LoginValid(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	// cookie := uuid.New().String()
	cookie := uuid.NewString()
	expiry := time.Now().Add(5 * time.Hour)
	// write cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   cookie,
		Expires: expiry,
	})
	session := model.Session{
		Username: creds.Username,
		Token:    cookie,
		Expiry:   expiry,
	}
	err = api.sessionsRepo.AddSessions(session)
	w.WriteHeader(200)
	api.dashboardView(w, r)

}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	//Read session_token and get Value:
	//sessionToken := "username" // TODO: replace this
	sessionToken := ""
	username := fmt.Sprintf("%s", r.Context().Value("username"))
	readSession, _ := api.sessionsRepo.ReadSessions()
	for _, v := range readSession {
		if username == v.Username {
			sessionToken = v.Token
		}
	}
	
	
	api.sessionsRepo.DeleteSessions(sessionToken)

	//Set Cookie name session_token value to empty and set expires time to Now:
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   "",
		Expires: time.Now(),
	})
	// TODO: answer here

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
	w.WriteHeader(http.StatusOK)
}
