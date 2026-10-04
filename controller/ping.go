package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func RunGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get request on RunGet")
	_, err := w.Write([]byte("hii"))
	if err != nil {
		fmt.Println(err)
	}
}

func RunPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post request on RunPost")
	id := chi.URLParam(r, "id")
	//fmt.Println(id)
	_, err := w.Write([]byte(id))
	if err != nil {
		panic("Error: " + err.Error())
	}
}
