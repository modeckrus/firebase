package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	firebase "firebase.google.com/go"
)

//AuthEvent is
type AuthEvent struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}
type User struct {
	UID      string `json:"uid"`      //User id
	Nick     string `json:"nickname"` //Nick name wich will display(for example in posts)
	Name     string `json:"name"`     //Actual Name
	Surname  string `json:"surname"`  //Actual Surname
	Avatar   string `json:"avatar"`   //Path in firebase storage to avatar // For displaying we will use an autogenereted thubnails
	Email    string `json:"email"`    //Email of user
	IsSetted bool   `json:"issetted"` //Is User was setted(specified name avatar and nick)
}

//SubModel is model representing subscriber of user
type SubModel struct {
	UID    string `json:"uid"`
	Nick   string `json:"nick"`
	Avatar string `json:"avatar"`
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

// AuthEventFunc is triggered by a change to a Firestore document.
func AuthEventFunc(ctx context.Context, e AuthEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%v", e)

	nick := strings.Split(e.Email, "@")[0]
	fstore.Collection("user").Doc(e.UID).Set(ctx, User{
		UID:      e.UID,
		Email:    e.Email,
		Nick:     nick,
		Name:     nick,
		Surname:  nick,
		IsSetted: false,
	})
	fstore.Collection("user").Doc(e.UID).Collection("subscribers").Doc(e.UID).Set(ctx, SubModel{
		UID:    e.UID,
		Nick:   nick,
		Avatar: "avatars/default.jpg",
	})
	return nil
}
