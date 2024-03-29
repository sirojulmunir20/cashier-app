package api

import (
	repo "a21hc3NpZ25tZW50/repository"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

type API struct {
	usersRepo    repo.UserRepository
	sessionsRepo repo.SessionsRepository
	products     repo.ProductRepository
	cartsRepo    repo.CartRepository
	mux          *http.ServeMux
}

type Page struct {
	File string
}

func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("views", p.File)
	// Fungsi `template.ParseFiles()`, digunakan untuk parsing file template
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Data map yang berisikan key `title` dan `name` yang akan di sisipkan ke dalam template yang sudah di parsing.
	// var data = map[string]interface{}{
	// 	"title": "Learning Golang Web",
	// 	"name":  "Student",
	// }

	// Method `Execute()` milik `*template.Template`, digunakan untuk menyisipkan data pada template, untuk kemudian ditampilkan ke browser. Data bisa disipkan dalam bentuk `struct`, `map`, atau `interface{}`.
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewAPI(usersRepo repo.UserRepository, sessionsRepo repo.SessionsRepository, products repo.ProductRepository, cartsRepo repo.CartRepository) API {
	mux := http.NewServeMux()
	api := API{
		usersRepo,
		sessionsRepo,
		products,
		cartsRepo,
		mux,
	}

	index := Page{File: "index.html"}
	mux.Handle("/", api.Get(index))

	register := Page{File: "register.html"}
	mux.Handle("/page/register", api.Get(register))
	login := Page{File: "login.html"}
	mux.Handle("/page/login", api.Get(login))
	// Please create routing for:
	// - Register page with endpoint `/page/register`, GET method and render `register.html` file on views folder
	// - Login page with endpoint `/page/login`, GET method and render `login.html` file on views folder
	// TODO: answer here

	mux.Handle("/user/register", api.Post(http.HandlerFunc(api.Register)))
	mux.Handle("/user/login", api.Post(http.HandlerFunc(api.Login)))
	mux.Handle("/user/logout", api.Get(api.Auth(http.HandlerFunc(api.Logout))))

	mux.Handle("/user/img/profile", api.Get(api.Auth(http.HandlerFunc(api.ImgProfileView))))
	mux.Handle("/user/img/update-profile", api.Post(api.Auth(http.HandlerFunc(api.ImgProfileUpdate))))

	// Please create routing for endpoint `/cart/add` with POST method, Authentication and handle api.AddCart
	// TODO: answer here
	mux.Handle("/cart/add", api.Post(api.Auth(http.HandlerFunc(api.AddCart))))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":8080", api.Handler())
}
