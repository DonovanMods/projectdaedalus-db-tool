package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/auth/credentials"
	gfs "cloud.google.com/go/firestore"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var fsClient *gfs.Client

func getClient() (*gfs.Client, error) {
	if fsClient != nil {
		return fsClient, nil
	}

	logger.Log.Info("Initializing Firestore client...")

	projectID := viper.GetString("firebase.credentials.project_id")
	credsJson, err := json.Marshal(viper.Get("firebase.credentials"))
	if err != nil {
		log.Panic(err)
	}

	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		CredentialsJSON: credsJson,
	})
	if err != nil {
		fmt.Println("creds: ", credsJson)
		log.Panic(err)
	}

	fsClient, err := gfs.NewClient(context.Background(), projectID, option.WithAuthCredentials(creds))
	if err != nil {
		log.Panic(err)
	}

	return fsClient, nil
}
