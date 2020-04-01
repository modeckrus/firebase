package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"socialapp/firebasestorage"

	"cloud.google.com/go/firestore"
	cstorage "cloud.google.com/go/storage"
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
var bucket *gstorage.BucketHandle
var auth *Auth.Client
var cstor *cstorage.Client

func init() {
	ctx := context.Background()
	conf := &firebase.Config{
		ProjectID:     "gvisionmodeck",
		DatabaseURL:   "https://gvisionmodeck.firebaseio.com/",
		StorageBucket: "gvisionmodeck.appspot.com",
	}
	opt := option.WithCredentialsFile("./secured/adminsdk.json") //Specify this file by ur adminsdk, u can find it in settigns of ur firebase project
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
	cstor, err = cstorage.NewClient(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}
	bucket = cstor.Bucket("gvisionmodeck.appspot.com")
}

func main() {
	b, err := firebasestorage.Read(cstor, bucket, "david.jpg")
	io := bytes.NewReader(b)
	src, err := imaging.Decode(io)
	err = imaging.Save(src, "images/david.jpg")
	if err != nil {
		log.Fatal(err)
	}
	var sizes = []int{
		100,
		200,
		300,
		500,
	}
	for i, size := range sizes {
		img := imaging.Thumbnail(src, size, size, imaging.Lanczos)
		err = imaging.Save(img, fmt.Sprintf("images/@thubnail_david_%v.jpg", i))
		if err != nil {
			log.Fatal(err)
		}
	}

}
