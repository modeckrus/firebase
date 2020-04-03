package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	firebase "firebase.google.com/go"
	usermodel "github.com/modeckrus/firebase/usermodel"
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

type AuthEvent struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log an interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     User      `json:"fields"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
}

// MyData represents a value from Firestore. The type definition depends on the
// format of your database.
type User struct {
	email string
	subsc []string
}
type Post struct {
	Body      string   `json:"body"`
	Dname     string   `json:"dname"`
	Title     string   `json:"title"`
	Likes     int      `json:"likes"`
	Hasattach bool     `json:"hasattach"`
	Images    []string `json:"images"`
}

// GOOGLE_CLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var fstore *firestore.Client

func init() {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	fstore, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

// HelloFirestore is triggered by a change to a Firestore document.
func AuthEventFunc(ctx context.Context, e AuthEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%v", e)

	nick := strings.Split(e.Email, "@")[0]
	fstore.Collection("user").Doc(e.UID).Set(ctx, usermodel.User{
		UID:   e.UID,
		Email: e.Email,
		Nick:  nick,
	})
	fstore.Collection("user").Doc(e.UID).Collection("subscribers").Doc(e.UID).Set(ctx, usermodel.SubModel{
		UID:    e.UID,
		Nick:   nick,
		Avatar: "avatars/default.jpg",
	})
	return nil
}
