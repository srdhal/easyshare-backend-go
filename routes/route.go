package routes

import (
	"github.com/go-chi/chi"
	"github.com/srdhal/easyshare-backend-go/controller"
)

func ImageRoute(r chi.Router){
	r.Post("/upload",controller.UploadImage)
	r.Get("/file/{fileid}",controller.GetImage)
}