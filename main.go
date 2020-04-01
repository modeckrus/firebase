package main

import (
	"context"
	"fmt"
	"log"
	"socialapp/testhandler"

	"cloud.google.com/go/firestore"
	cstorage "cloud.google.com/go/storage"
	gstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	Auth "firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	firestorage "firebase.google.com/go/storage"
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

type ThubHand struct {
	Path  string `json:"path"`
	Sizes []struct {
		Height string `json:"height"`
		Width  int    `json:"width"`
	} `json:"sizes"`
}

func main() {
	ctx := context.Background()
	doc := fstore.Collection("service").Doc("service").Collection("thubnail").Doc("x9f09yczE0fGiNyehsul")
	snap, err := doc.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var thub ThubHand
	err = snap.DataTo(&thub)
	if err != nil {
		log.Fatal(err)
	}
	fstore.Collection("service").Doc("service").Collection("thubnail").Add(ctx, thub)
	fmt.Println(thub)
	/*
		b, err := firebasestorage.Read(cstor, bucket, "david.jpg")
		io := bytes.NewReader(b)
		src, err := imaging.Decode(io)
		err = imaging.Save(src, "images/david.jpg")
		if err != nil {
			log.Fatal(err)
		}

		var sizes []thubnails.Size
		sizes = []thubnails.Size{
			thubnails.Size{
				Width:  100,
				Height: 100,
			},
			thubnails.Size{
				Width:  200,
				Height: 200,
			},
			thubnails.Size{
				Width:  300,
				Height: 300,
			},
			thubnails.Size{
				Width:  500,
				Height: 500,
			},
		}
		thubnails.Thubnails(cstor, bucket, "david.jpg", sizes)
	*/
	testhandler.Thubnails()
}
