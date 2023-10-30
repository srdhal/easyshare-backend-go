package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var FileCollection *mongo.Collection
func DBConnection(){
	dbUrl:=os.Getenv("MONGO_URL")

	clientOptions:=options.Client().ApplyURI(dbUrl)
	client,err:=mongo.Connect(context.Background(),clientOptions)
	
	if(err!=nil){
		log.Fatal("db connection error: ",err)
	}
	
	err=client.Ping(context.Background(),nil)
	
	if(err!=nil){
		log.Fatal("db connection error: ",err)
	}else{
		log.Printf("db connected")
	}

	FileCollection=client.Database("easyshare").Collection("files")
	
}


