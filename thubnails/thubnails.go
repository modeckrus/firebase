package thubnails

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	cstorage "cloud.google.com/go/storage"
	gstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	Auth "firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	firestorage "firebase.google.com/go/storage"
	"github.com/disintegration/imaging"
	"github.com/modeckrus/firebase/firebasestorage"
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

var client *db.Client
var fstore *firestore.Client
var storage *firestorage.Client
var bucket *gstorage.BucketHandle
var auth *Auth.Client
var cstor *cstorage.Client

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	if projectID == "" {
		projectID = "gvisionmodeck"
	}
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
	bucket = cstor.Bucket(fmt.Sprintf("%v.appspot.com", projectID))
}

//ThubHandStruct event
type ThubHandStruct struct {
	Path struct {
		StringValue string `json:"stringValue"`
	} `json:"Path"`
	Ready struct {
		BooleanValue bool `json:"booleanValue"`
	} `json:"Ready"`
	Sizes struct {
		ArrayValue struct {
			Values []struct {
				MapValue struct {
					Fields struct {
						Height struct {
							IntegerValue string `json:"integerValue"`
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

//Thubnail reference of firestore
type Thubnail struct {
	Path  string `json:"path"`
	Sizes []Size `json:"sizes"`
	Ready bool   `json:"ready"`
}

func convertToThubStrcut(hand ThubHandStruct) (Thubnail, error) {
	var thub Thubnail
	thub.Path = hand.Path.StringValue
	for _, size := range hand.Sizes.ArrayValue.Values {
		h, err := strconv.Atoi(size.MapValue.Fields.Height.IntegerValue)
		if err != nil {
			return thub, err
		}
		w, err := strconv.Atoi(size.MapValue.Fields.Width.IntegerValue)
		if err != nil {
			return thub, err
		}
		thub.Sizes = append(thub.Sizes, Size{
			Height: h,
			Width:  w,
		})
	}
	thub.Ready = false
	return thub, nil
}

//ThumbHandler handle the event and convert it
func ThumbHandler(ctx context.Context, e FirestoreEvent) error {
	fullPath := strings.Split(e.Value.Name, "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	log.Println(pathParts)
	js, err := json.Marshal(e.Value.Fields)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("event fields: %v", string(js))

	thub, err := convertToThubStrcut(e.Value.Fields)
	if err != nil {
		log.Fatal(err)
		return err
	}
	Thubnails(cstor, bucket, thub.Path, thub.Sizes)
	thub.Ready = true
	log.Printf("Thub is: %v", thub)
	fstore.Collection("service").Doc("service").Collection("thubnail").Doc(pathParts[len(pathParts)-1]).Set(
		ctx,
		thub,
	)

	return nil
}

//Thubnails make thunail of images in firebase storage
func Thubnails(cstor *gstorage.Client, bucket *gstorage.BucketHandle, filename string, sizes []Size) {
	b, err := firebasestorage.Read(cstor, bucket, filename)
	if err != nil {
		log.Fatal(err)
	}
	io := bytes.NewReader(b)
	src, err := imaging.Decode(io)
	path := strings.Split(filename, "/")
	f := path[len(path)-1] //Filename
	path = path[:len(path)-1]
	p := strings.Join(path, "/") //Path to file
	fmt.Printf("%v/%v", p, f)
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
		err = firebasestorage.Write(cstor, bucket, fmt.Sprintf("%v/@thub_%vX%v_%v", p, size.Width, size.Height, f), &buf)
		if err != nil {
			log.Println(err)
		}
		//imaging.Encode( )
	}
}
