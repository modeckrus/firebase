package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	gstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	Auth "firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	firestorage "firebase.google.com/go/storage"
	"github.com/disintegration/imaging"
	"google.golang.org/api/option"
)

var client *db.Client
var fstore *firestore.Client
var storage *firestorage.Client
var buckit *gstorage.BucketHandle
var auth *Auth.Client

func init() {
	ctx := context.Background()
	conf := &firebase.Config{
		ProjectID:     "gvisionmodeck",
		DatabaseURL:   "https://gvisionmodeck.firebaseio.com/",
		StorageBucket: "gvisionmodeck.appspot.com",
	}
	opt := option.WithCredentialsFile("./adminsdk.json")
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatal(err)
	}
	client, err = app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fstore, err = app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	auth, err = app.Auth(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	srcimg, err := imaging.Open("./1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	img := imaging.Thumbnail(srcimg, 300, 300, imaging.Lanczos)
	err = imaging.Save(img, fmt.Sprintf("@thubnail_%v", "1.jpg"))
	if err != nil {
		log.Fatal(err)
	}

}
