package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	// View with response image `img-avatar.png` from path `assets/images`
	// TODO: answer here
	fileBytes, _ := ioutil.ReadFile("./assets/images/img-avatar.png")
	w.WriteHeader(200)
	w.Write(fileBytes)
}

func (api *API) ImgProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// Update image `img-avatar.png` from path `assets/images`
	// TODO: answer here
	//update file func diatasnya
	r.ParseMultipartForm(3 << 20)
	uploadFile, _, err := r.FormFile("file-avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	defer uploadFile.Close()
	
	w.WriteHeader(200)

// foleder yang dituju
targetFile, _ := os.OpenFile("./assets/images/img-avatar.png", os.O_WRONLY|os.O_CREATE, 0666)

defer targetFile.Close()

//replace image dari client
_, err = io.Copy(targetFile, uploadFile)
if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

api.dashboardView(w, r)
}
