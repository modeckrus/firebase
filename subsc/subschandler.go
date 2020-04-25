package subsc

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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
	Fields     interface{} `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

// MyData represents a value from Firestore. The type definition depends on the
// format of your database.

//PostE is post event that trigger this func
type SubscE struct {
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
func SubscCreated(ctx context.Context, e FirestoreEvent) error {
	fullPath := strings.Split(e.Value.Name, "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	//collection := pathParts[0]
	uid := strings.Join(pathParts[1:2], "/")
	subscId := strings.Join(pathParts[3:], "/")
	log.Print("uid: ", uid)
	log.Print("subscID", subscId)
	json, err := json.Marshal(e.Value.Fields)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(json))

	return nil
}
