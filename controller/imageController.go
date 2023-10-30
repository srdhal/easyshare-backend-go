package controller

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/srdhal/easyshare-backend-go/database"
	"github.com/srdhal/easyshare-backend-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func UploadImage(w http.ResponseWriter,r *http.Request) {
     file,header,err:=r.FormFile("file")
	 if err!=nil{
		log.Fatal("file not found: ",err)
	 }
	 defer file.Close()

	 dst,err:=os.Create(filepath.Join("./uploads/",header.Filename))
	 if err!=nil{
		log.Fatal("file upload error: ",err)
	 }
	 defer dst.Close()

	 _,_=io.Copy(dst,file)

	 uploadFile:=models.File{Path: "uploads/"+header.Filename,Name: header.Filename,DownloadCounter: 0}
	 uploadedFile,err:=database.FileCollection.InsertOne(context.Background(),uploadFile)

     if err!=nil{
		log.Println("file not inserted")
	 } 

	 path:=os.Getenv("BACKEND_URL")+"/file/"+uploadedFile.InsertedID.(primitive.ObjectID).Hex()

     response:=struct{
		Path string `json:"path"`
	 }{
		Path: path,
	 }
     w.Header().Set("Content-Type","application/json")
	 w.WriteHeader(http.StatusOK)

	 if err:=json.NewEncoder(w).Encode(response); err!=nil{
		http.Error(w,"JSON encoding error",http.StatusInternalServerError)
		return
	 }

	 log.Println("upload success: ",uploadedFile)
}

func GetImage(w http.ResponseWriter,r *http.Request){
    // log.Println("getimage")
	fileid:=chi.URLParam(r, "fileid")
	fileID,_:=primitive.ObjectIDFromHex(fileid)
    filter:=bson.M{"_id": fileID}
	var retrivedFile bson.M
	err:=database.FileCollection.FindOne(context.Background(),filter).Decode(&retrivedFile)
	if err!=nil{
		log.Println(err)
		http.Error(w,"file retriving error",http.StatusInternalServerError)
		return
	}

	// log.Println("retrived file: ",retrivedFile)

    path:=retrivedFile["path"].(string)
	update:=bson.M{"$inc":bson.M{"downloadCounter": 1}}
	file,err:=os.Open(path)
    if err!=nil{
		http.Error(w,"file not found",http.StatusNotFound)
		return
	}
    defer file.Close()
	w.Header().Set("Content-Disposition","attachment; filename="+retrivedFile["name"].(string))    
    w.Header().Set("Content-Type","application/pdf")
	
	_,err=io.Copy(w,file)
	if err!=nil{
		http.Error(w,"error serving the file",http.StatusInternalServerError)
		return
	}
	_,err=database.FileCollection.UpdateOne(context.Background(),filter,update)
	if err!=nil{
		http.Error(w,"error in update counter",http.StatusInternalServerError)
		return
	}
}

