package firestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/auth/credentials"
	gfs "cloud.google.com/go/firestore"
	urlverifier "github.com/davidmytton/url-verifier"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/mod"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type SoM interface {
	string | mod.Mod
}

type DBList[SM SoM] interface {
	fmt.Stringer
	json.Marshaler
	Add(item SM) error
	Commit() (*gfs.WriteResult, error)
	Count() int
	Fetch() error
	Parse(string) error
	Remove(item SM) error
}

var (
	fsClient              *gfs.Client
	dbMetaCollectionTypes = map[string]DBList[string]{
		"repositories": &repoCache,
		"modinfo":      &modInfoCache,
		"toolinfo":     &toolInfoCache,
	}
)

type ErrConfigNotFound struct {
	item string
}

func (e ErrConfigNotFound) Error() string {
	return fmt.Sprintf("config item %q not found", e.item)
}

func CommitAll() error {
	if fsClient == nil {
		logger.Panic(errors.New("firestore client not initialized"))
	}

	logger.Info("committing All Firestore changes")

	// Commit Meta collections
	for _, collection := range dbMetaCollectionTypes {
		if _, err := collection.Commit(); err != nil {
			return err
		}
	}

	// Commit mod collections
	if _, err := modCache.Commit(); err != nil {
		return err
	}

	return fsClient.Close()
}

func FetchAll() error {
	// Fetch Meta lists
	for _, collection := range dbMetaCollectionTypes {
		if err := collection.Fetch(); err != nil {
			return err
		}
	}

	// Fetch Mods
	if err := modCache.Fetch(); err != nil {
		return err
	}

	return nil
}

/*
// Private functions
*/
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

	logger.Info(fmt.Sprintf("fetching document from %q", collection))

	return client.Doc(collection).Get(context.Background())
}

func getDocuments(collectionString string) (*gfs.DocumentIterator, error) {
	collection, err := getCollection(collectionString)
	if err != nil {
		return nil, err
	}

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("fetching documents from %q", collection))

	return client.Collection(collection).Documents(context.Background()), nil
}

func verifyURL(url string) error {
	verifier := urlverifier.NewVerifier()
	verifier.EnableHTTPCheck()

	logger.Info(fmt.Sprintf("validating %q", url))

	ret, err := verifier.Verify(url)
	if err != nil {
		return err
	}

	if !ret.HTTP.IsSuccess {
		return fmt.Errorf("could not reach %q; code: %d", url, ret.HTTP.StatusCode)
	}

	return nil
}
