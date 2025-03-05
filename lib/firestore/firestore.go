package firestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/auth/credentials"
	gfs "cloud.google.com/go/firestore"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	ErrDuplicate      = errors.New("item already exists")
	fsClient          *gfs.Client
	dbCollectionTypes = map[string]*metaList{
		"repositories": &repo,
		"modinfo":      &modInfo,
		"toolinfo":     &toolInfo,
	}
)

type ErrConfigNotFound struct {
	item string
}

func (e ErrConfigNotFound) Error() string {
	return fmt.Sprintf("config item %q not found", e.item)
}

func Commit() error {
	if fsClient == nil {
		logger.Panic(errors.New("firestore client not initialized"))
	}

	logger.Info("committing All Firestore changes")

	for _, collection := range dbCollectionTypes {
		if _, err := collection.Commit(); err != nil {
			return err
		}
	}

	return fsClient.Close()
}

func getClient() (*gfs.Client, error) {
	if fsClient != nil {
		return fsClient, nil
	}

	logger.Info("initializing Firestore client")

	projectID := viper.GetString("firebase.credentials.project_id")
	credsJson, err := json.Marshal(viper.Get("firebase.credentials"))
	if err != nil {
		logger.Panic(err)
	}

	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		CredentialsJSON: credsJson,
	})
	if err != nil {
		fmt.Println("creds: ", credsJson)
		logger.Panic(err)
	}

	fsClient, err = gfs.NewClient(context.Background(), projectID, option.WithAuthCredentials(creds))
	if err != nil {
		logger.Panic(err)
	}

	return fsClient, nil
}

func getCollection(collectionString string) (string, error) {
	collection := viper.GetString(collectionString)

	if collection == "" {
		return "", &ErrConfigNotFound{item: collectionString}
	}

	return collection, nil
}

func getDocument(collectionString string) (*gfs.DocumentSnapshot, error) {
	collection, err := getCollection(collectionString)
	if err != nil {
		return nil, err
	}

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("fetching documents from %q", collection))

	return client.Doc(collection).Get(context.Background())
}
