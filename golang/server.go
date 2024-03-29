package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/http2"
)

var (
	bucketName         string = ""
	serviceAccountName string = ""
	objectName         string = "file.txt"
)

func fronthandler(w http.ResponseWriter, r *http.Request) {

	expires := time.Now().Add(time.Minute * 10)

	// this will work
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting client %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	s, err := storageClient.Bucket(bucketName).SignedURL(objectName, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Expires: expires,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting signedURL %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// this will also work but is not the right way to do this...
	// ctx := context.Background()

	// rootTokenSource, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/iam")
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }
	// delegates := []string{}

	// s, err := storage.SignedURL(bucketName, objectName, &storage.SignedURLOptions{
	// 	Scheme:         storage.SigningSchemeV4,
	// 	GoogleAccessID: serviceAccountName,
	// 	SignBytes: func(b []byte) ([]byte, error) {
	// 		client := oauth2.NewClient(context.TODO(), rootTokenSource)
	// 		service, err := iamcredentials.New(client)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("oauth2/google: Error creating IAMCredentials: %v", err)
	// 		}
	// 		signRequest := &iamcredentials.SignBlobRequest{
	// 			Payload:   base64.StdEncoding.EncodeToString(b),
	// 			Delegates: delegates,
	// 		}
	// 		name := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountName)
	// 		at, err := service.Projects.ServiceAccounts.SignBlob(name, signRequest).Do()
	// 		if err != nil {
	// 			return nil, fmt.Errorf("oauth2/google: Error calling iamcredentials.SignBlob: %v", err)
	// 		}
	// 		sDec, err := base64.StdEncoding.DecodeString(at.SignedBlob)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("oauth2/google: Error decoding iamcredentials.SignBlob response: %v", err)
	// 		}
	// 		return sDec, nil
	// 	},
	// 	Method:  http.MethodGet,
	// 	Expires: expires,
	// })
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }
	fmt.Println(s)
	w.Header().Set("context-type", "text/plain")
	fmt.Fprint(w, s)
}

func main() {

	http.HandleFunc("/", fronthandler)

	bucketName = os.Getenv("BUCKET_NAME")
	serviceAccountName = os.Getenv("SA_EMAIL")

	var server *http.Server
	server = &http.Server{
		Addr: ":8080",
	}
	http2.ConfigureServer(server, &http2.Server{})
	fmt.Printf("Starting Server..\n")
	err := server.ListenAndServe()
	fmt.Printf("Unable to start Server %v", err)
}
