package posthandler

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	postedition "github.com/modeckrus/firebase/postedition"
	usermodel "github.com/modeckrus/firebase/usermodel"
	"google.golang.org/api/iterator"
)

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
	Fields     PostE     `json:"fields"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
}

// MyData represents a value from Firestore. The type definition depends on the
// format of your database.

//PostE is post event that trigger this func
type PostE struct {
	Body struct {
		StringValue string `json:"stringValue"`
	} `json:"Body"`
	Hasattach struct {
		BooleanValue bool `json:"booleanValue"`
	} `json:"Hasattach"`
	Images struct {
		ArrayValue struct {
			Values []struct {
				StringValue string `json:"stringValue"`
			} `json:"values"`
		} `json:"arrayValue"`
	} `json:"Images"`
	Title struct {
		StringValue string `json:"stringValue"`
	} `json:"Title"`
	Subtitle struct {
		StringValue string `json:"stringValue"`
	} `json:"Subtitle"`
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

// PostCreated is triggered by a change to a Firestore document. It updates
// the `original` value of the document to upper case.
func PostCreated(ctx context.Context, e FirestoreEvent) error {
	fullPath := strings.Split(e.Value.Name, "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	//collection := pathParts[0]
	uid := strings.Join(pathParts[1:2], "/")
	postid := strings.Join(pathParts[3:], "/")
	log.Print("uid: ", uid)
	log.Print("postid: ", postid)
	log.Printf("event fields: %v", e.Value.Fields)

	uDoc, err := fstore.Collection("user").Doc(uid).Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	user := &usermodel.User{}
	err = uDoc.DataTo(user)
	if err != nil {
		log.Fatal(err)
	}
	var images []string
	for i := 0; i < len(e.Value.Fields.Images.ArrayValue.Values); i++ {
		image := e.Value.Fields.Images.ArrayValue.Values[i]
		images = append(images, image.StringValue)
	}
	if err != nil {
		log.Printf("error converting lieks: %v", err)
	}
	postE := postedition.PostEdition{
		Body:      e.Value.Fields.Body.StringValue,
		Title:     e.Value.Fields.Title.StringValue,
		Subtitle:  e.Value.Fields.Subtitle.StringValue,
		Hasattach: e.Value.Fields.Hasattach.BooleanValue,

		Images: images,
	}
	log.Println(postE)
	postPub := CretePostPubfromPostEdition(postE, *user)
	_, _, err = fstore.Collection("posts").Doc(user.UID).Collection("post").Add(ctx, postPub)
	if err != nil {
		log.Println("error while adding document in firebase")
		log.Println(err)
		return err
	}
	fstore.Collection("feedline").Doc(uid).Collection("post").Doc(pathParts[3]).Set(ctx, postPub)
	subsciter := fstore.Collection("user").Doc(user.UID).Collection("subscribers").Documents(ctx)
	for {
		subDoc, err := subsciter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Error while iterating sub: ", err)
			break
		}
		sub := &usermodel.SubModel{}
		err = subDoc.DataTo(sub)
		if err != nil {
			log.Println("Error while marshaling request into sub: ", err)
			break
		}
		suid := sub.UID
		log.Printf("Subscriber: %v", suid)
		log.Printf("Documents: %v", e.Value.Fields)
		fstore.Collection("feedline").Doc(suid).Collection("post").Doc(pathParts[3]).Set(ctx, postPub)
	}

	return nil
}

//CretePostPubfromPostEdition creating post pub from post editon
func CretePostPubfromPostEdition(postE postedition.PostEdition, user usermodel.User) *postedition.PostPub {
	postPub := &postedition.PostPub{}
	var nick string
	var avatar string

	if user.Nick == "" {
		nick = user.Email
	} else {
		nick = user.Nick
	}
	if user.Avatar == "" {
		avatar = "avatars/default.jpg"
	} else {
		avatar = user.Avatar
	}
	postPub = &postedition.PostPub{
		Title:     postE.Title,
		Subtitle:  postE.Subtitle,
		Body:      postE.Body,
		Nick:      nick,
		Avatar:    avatar,
		Likes:     0,
		UserID:    user.UID,
		Hasattach: postE.Hasattach,
		Images:    postE.Images,
		CreatedAt: time.Now(),
	}
	return postPub
}
