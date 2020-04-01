package firebasestorage

///Interface for firebase storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func main() {
	/*
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
			os.Exit(1)
		}
		var o string
		flag.StringVar(&o, "o", "", "source object; in the format of <bucket:object>")
		flag.Parse()

		names := strings.Split(o, ":")
		if len(names) < 2 {
			usage("missing -o flag")
		}
		bucket, object := names[0], names[1]

		if len(os.Args) < 3 {
			usage("missing subcommand")
		}

		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}

		switch os.Args[2] {
		case "write":
			if err := WriteFile(client, bucket, object, "notes.txt"); err != nil {
				log.Fatalf("Cannot write object: %v", err)
			}
		case "read":
			data, err := Read(client, bucket, object)
			if err != nil {
				log.Fatalf("Cannot read object: %v", err)
			}
			fmt.Printf("Object contents: %s\n", data)
		case "metadata":
			attrs, err := Attrs(client, bucket, object)
			if err != nil {
				log.Fatalf("Cannot get object metadata: %v", err)
			}
			fmt.Printf("Object metadata: %v\n", attrs)
		case "makepublic":
			if err := MakePublic(client, bucket, object); err != nil {
				log.Fatalf("Cannot to make object public: %v", err)
			}
		case "delete":
			if err := Delete(client, bucket, object); err != nil {
				log.Fatalf("Cannot to delete object: %v", err)
			}
		}
	*/
}

//WriteFile upload the file to firebase storage
func WriteFile(client *storage.Client, bucket *storage.BucketHandle, object, filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	Write(client, bucket, object, f)

	return nil
}

//Write upload file using io.Reader, u can put a file here
func Write(client *storage.Client, bucket *storage.BucketHandle, object string, f io.Reader) error {
	// [START upload_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := bucket.Object(object).NewWriter(ctx)
	var err error
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END upload_file]
	return nil
}

//List write all objects in bucket and put it in io.Writer
func List(w io.Writer, client *storage.Client, bucket string) error {
	// [START storage_list_files]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files]
	return nil
}

//ListByPrefix with prefix
// Prefixes and delimiters can be used to emulate directory listings.
// Prefixes can be used filter objects starting with prefix.
// The delimiter argument can be used to restrict the results to only the
// objects in the given "directory". Without the delimiter, the entire  tree
// under the prefix is returned.
//
// For example, given these blobs:
//   /a/1.txt
//   /a/b/2.txt
//
// If you just specify prefix="a/", you'll get back:
//   /a/1.txt
//   /a/b/2.txt
//
// However, if you specify prefix="a/" and delim="/", you'll get back:
//   /a/1.txt
func ListByPrefix(w io.Writer, client *storage.Client, bucket, prefix, delim string) error {
	// [START storage_list_files_with_prefix]
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimiter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files_with_prefix]
	return nil
}

//Read ...
func Read(client *storage.Client, bucket *storage.BucketHandle, object string) ([]byte, error) {
	// [START download_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := bucket.Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

//ReadIo return io.Reader instead bytes
func ReadIo(client *storage.Client, bucket *storage.BucketHandle, object string) (io.Reader, error) {
	// [START download_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := bucket.Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return rc, nil
}

//Attrs ...
func Attrs(client *storage.Client, bucket, object string) (*storage.ObjectAttrs, error) {
	// [START get_metadata]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Bucket: %v\n", attrs.Bucket)
	log.Printf("CacheControl: %v\n", attrs.CacheControl)
	log.Printf("ContentDisposition: %v\n", attrs.ContentDisposition)
	log.Printf("ContentEncoding: %v\n", attrs.ContentEncoding)
	log.Printf("ContentLanguage: %v\n", attrs.ContentLanguage)
	log.Printf("ContentType: %v\n", attrs.ContentType)
	log.Printf("Crc32c: %v\n", attrs.CRC32C)
	log.Printf("Generation: %v\n", attrs.Generation)
	log.Printf("KmsKeyName: %v\n", attrs.KMSKeyName)
	log.Printf("Md5Hash: %v\n", attrs.MD5)
	log.Printf("MediaLink: %v\n", attrs.MediaLink)
	log.Printf("Metageneration: %v\n", attrs.Metageneration)
	log.Printf("Name: %v\n", attrs.Name)
	log.Printf("Size: %v\n", attrs.Size)
	log.Printf("StorageClass: %v\n", attrs.StorageClass)
	log.Printf("TimeCreated: %v\n", attrs.Created)
	log.Printf("Updated: %v\n", attrs.Updated)
	log.Printf("Event-based hold enabled? %t\n", attrs.EventBasedHold)
	log.Printf("Temporary hold enabled? %t\n", attrs.TemporaryHold)
	log.Printf("Retention expiration time %v\n", attrs.RetentionExpirationTime)
	log.Print("\n\nMetadata\n")
	for key, value := range attrs.Metadata {
		log.Printf("\t%v = %v\n", key, value)
	}

	return attrs, nil
	// [END get_metadata]
}

//SetEventBasedHold ...
func SetEventBasedHold(client *storage.Client, bucket, object string) error {
	// [START storage_set_event_based_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		EventBasedHold: true,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_set_event_based_hold]
	return nil
}

//ReleaseEventBasedHold ...
func ReleaseEventBasedHold(client *storage.Client, bucket, object string) error {
	// [START storage_release_event_based_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		EventBasedHold: false,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_release_event_based_hold]
	return nil
}

//SetTemporaryHold ...
func SetTemporaryHold(client *storage.Client, bucket, object string) error {
	// [START storage_set_temporary_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		TemporaryHold: true,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_set_temporary_hold]
	return nil
}

//ReleaseTemporaryHold ...
func ReleaseTemporaryHold(client *storage.Client, bucket, object string) error {
	// [START storage_release_temporary_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		TemporaryHold: false,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_release_temporary_hold]
	return nil
}

//MakePublic ...
func MakePublic(client *storage.Client, bucket, object string) error {
	// [START public]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	// [END public]
	return nil
}

//Move ...
func Move(client *storage.Client, bucket, object string) error {
	// [START move_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	dstName := object + "-rename"

	src := client.Bucket(bucket).Object(object)
	dst := client.Bucket(bucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}
	// [END move_file]
	return nil
}

//CopyToBucket ...
func CopyToBucket(client *storage.Client, dstBucket, srcBucket, srcObject string) error {
	// [START copy_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	dstObject := srcObject + "-copy"
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	// [END copy_file]
	return nil
}

//Delete ...
func Delete(client *storage.Client, bucket, object string) error {
	// [START delete_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END delete_file]
	return nil
}

// WriteEncryptedObject writes an object encrypted with user-provided AES key to a bucket.
func WriteEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) error {
	// [START storage_upload_encrypted_file]
	ctx := context.Background()
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Encrypt the object's contents.
	wc := obj.Key(secretKey).NewWriter(ctx)
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_encrypted_file]
	return nil
}

// WriteWithKMSKey writes an object encrypted with KMS-provided key to a bucket.
func WriteWithKMSKey(client *storage.Client, bucket, object string, keyName string) error {
	// [START storage_upload_with_kms_key]
	ctx := context.Background()
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Encrypt the object's contents
	wc := obj.NewWriter(ctx)
	wc.KMSKeyName = keyName
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_with_kms_key]
	return nil
}

//ReadEncryptedObject
func ReadEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) ([]byte, error) {
	// [START storage_download_encrypted_file]
	ctx := context.Background()
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	// [END storage_download_encrypted_file]
	return data, nil
}

//RotateEncryptionKey
func RotateEncryptionKey(client *storage.Client, bucket, object string, key, newKey []byte) error {
	// [START storage_rotate_encryption_key]
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	// obj is encrypted with key, we are encrypting it with the newKey.
	_, err = obj.Key(newKey).CopierFrom(obj.Key(key)).Run(ctx)
	if err != nil {
		return err
	}
	// [END storage_rotate_encryption_key]
	return nil
}

//DownloadUsingRequesterPays
func DownloadUsingRequesterPays(client *storage.Client, object, bucketName, localpath, billingProjectID string) error {
	// [START storage_download_file_requester_pays]
	ctx := context.Background()

	bucket := client.Bucket(bucketName).UserProject(billingProjectID)
	src := bucket.Object(object)

	f, err := os.OpenFile(localpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := src.NewReader(ctx)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, rc); err != nil {
		return err
	}
	if err := rc.Close(); err != nil {
		return err
	}
	fmt.Printf("Downloaded using %v as billing project.\n", billingProjectID)
	// [END storage_download_file_requester_pays]
	return nil
}

//GenerateV4GetObjectSignedURL
func GenerateV4GetObjectSignedURL(w io.Writer, client *storage.Client, bucketName, objectName, serviceAccount string) (string, error) {
	// [START storage_generate_signed_url_v4]
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("Unable to generate a signed URL: %v", err)
	}

	fmt.Fprintln(w, "Generated GET signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl %q\n", u)
	// [END storage_generate_signed_url_v4]
	return u, nil
}

func GenerateV4PutObjectSignedURL(w io.Writer, client *storage.Client, bucketName, objectName, serviceAccount string) (string, error) {
	// [START storage_generate_upload_signed_url_v4]
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			"Content-Type:application/octet-stream",
		},
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("Unable to generate a signed URL: %v", err)
	}

	fmt.Fprintln(w, "Generated PUT signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl -X PUT -H 'Content-Type: application/octet-stream' --upload-file my-file %q\n", u)
	// [END storage_generate_upload_signed_url_v4]
	return u, nil
}

// TODO(jbd): Add test for downloadUsingRequesterPays.

const helptext = `usage: objects -o=bucket:name [subcommand] <args...>
subcommands:
	- write
	- read
	- metadata
	- makepublic
	- delete
`

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, helptext)
	os.Exit(2)
}
