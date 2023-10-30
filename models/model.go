package models


type File struct{
	Path            string `json:"path" bson:"path"`
	Name            string `json:"name" bson:"name"`
	DownloadCounter int64 `json:"downloadCounter" bson:"downloadCounter"`
}