package testhandler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/modeckrus/firebase/firebasestorage"

	"cloud.google.com/go/firestore"
	cstorage "cloud.google.com/go/storage"
	gstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	Auth "firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	firestorage "firebase.google.com/go/storage"
	"github.com/disintegration/imaging"
)

//Size size
type Size struct {
	Width  int
	Height int
}

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log an interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     ThubHandStruct `json:"fields"`
	Name       string         `json:"name"`
	UpdateTime time.Time      `json:"updateTime"`
}

type ThumbVal struct {
}

var client *db.Client
var fstore *firestore.Client
var storage *firestorage.Client
var bucket *gstorage.BucketHandle
var auth *Auth.Client
var cstor *cstorage.Client

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	ctx := context.Background()
	conf := &firebase.Config{
		ProjectID:     projectID,
		DatabaseURL:   fmt.Sprintf("https://%v.firebaseio.com/", projectID),
		StorageBucket: fmt.Sprintf("%v.appspot.com", projectID),
	}
	//opt := option.WithCredentialsFile("./secured/adminsdk.json") //Specify this file by ur adminsdk, u can find it in settigns of ur firebase project
	app, err := firebase.NewApp(ctx, conf)
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

	cstor, err = cstorage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	bucket = cstor.Bucket("gvisionmodeck.appspot.com")
}

type ThubHandStruct struct {
	Path struct {
		StringValue string `json:"stringValue"`
	} `json:"Path"`
	Sizes struct {
		ArrayValue struct {
			Values []struct {
				MapValue struct {
					Fields struct {
						Height struct {
							StringValue string `json:"stringValue"`
						} `json:"Height"`
						Width struct {
							IntegerValue string `json:"integerValue"`
						} `json:"Width"`
					} `json:"fields"`
				} `json:"mapValue"`
			} `json:"values"`
		} `json:"arrayValue"`
	} `json:"Sizes"`
}

type Thubnail struct {
	Path  string `json:"path"`
	Sizes []struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"sizes"`
}

func convertToThubStrcut(hand ThubHandStruct) (Thubnail, error) {
	var thub Thubnail
	thub.Path = hand.Path.StringValue
	for _, size := range hand.Sizes.ArrayValue.Values {
		h, err := strconv.Atoi(size.MapValue.Fields.Height.StringValue)
		if err != nil {
			return thub, err
		}
		w, err := strconv.Atoi(size.MapValue.Fields.Width.IntegerValue)
		if err != nil {
			return thub, err
		}
		thub.Sizes = append(thub.Sizes, struct {
			Height int "json:\"height\""
			Width  int "json:\"width\""
		}{
			Height: h,
			Width:  w,
		})
	}
	return thub, nil
}

//Thubnails make thunail of images in firebase storage
func Thubnails(cstor *gstorage.Client, bucket *gstorage.BucketHandle, filename string, sizes []Size) {
	b, err := firebasestorage.Read(cstor, bucket, filename)
	io := bytes.NewReader(b)
	src, err := imaging.Decode(io)
	err = imaging.Save(src, filename)
	if err != nil {
		log.Fatal(err)
	}
	for _, size := range sizes {
		img := imaging.Thumbnail(src, size.Width, size.Height, imaging.Lanczos)
		/*
			err = imaging.Save(img, fmt.Sprintf("@thubnail_%vX%v_%v", size.Width, size.Height, filename))
			if err != nil {
				log.Fatal(err)
			}
		*/
		filetype, err := imaging.FormatFromFilename(filename)
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		imaging.Encode(&buf, img, filetype)
		//f, err := os.Open(fmt.Sprintf("@thub_%vX%v%v", size.Width, size.Height, filename))
		err = firebasestorage.Write(cstor, bucket, fmt.Sprintf("@thub_%vX%v%v", size.Width, size.Height, filename), &buf)
		if err != nil {
			log.Println(err)
		}
		//imaging.Encode( )
	}
}
