package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/srdhal/easyshare-backend-go/database"
	"github.com/srdhal/easyshare-backend-go/routes"
)

func main() {
	godotenv.Load(".env")

	port:=os.Getenv("PORT")

	if(port==""){
		log.Fatal("port value is not set")
	}
    
	router:=chi.NewRouter()
    router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))
	
	server:=&http.Server{
		Handler: router,
		Addr: ":"+port,
	}
    routes.ImageRoute(router)
	log.Printf("server is running on port %v",port)
	database.DBConnection()
	err:=server.ListenAndServe()
	if err!=nil{
		log.Fatal("server error: ",err)
	}
}

