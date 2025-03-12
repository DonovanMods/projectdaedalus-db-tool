package firestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gfs "cloud.google.com/go/firestore"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"golang.org/x/exp/slices"
)

var ErrDuplicate = errors.New("item already exists")

const metaCollectionBase = "firebase.collections.meta"

// Data Caches
var (
	repoCache     = metaList{name: "repositories"}
	modInfoCache  = metaList{name: "modinfo"}
	toolInfoCache = metaList{name: "toolinfo"}
)

type metaList struct {
	Items []string `firestore:"list"`
	name  string   `firestore:"-" json:"-"`
	dirty bool
}

func (m *metaList) Fetch() error {
	if m.Items != nil {
		logger.Info(fmt.Sprintf("Using cached data for %s", m.name))
		return nil
	}

	docSnap, err := getDocument(m.configCollectionString())
	if err != nil {
		return fmt.Errorf("getDocument: %w", err)
	}

	if !docSnap.Exists() {
		return errors.New("document does not exist")
	}

	if err := docSnap.DataTo(&m); err != nil {
		return fmt.Errorf("DataTo: %w", err)
	}

	logger.Info(fmt.Sprintf("successfully retrieved %s list", m.name))

	return nil
}

// Add adds or updates an item in the list
// If item is already in the list, it will be removed and replaced
func (m *metaList) Add(item string) error {
	if item == "" {
		return errors.New("item cannot be blank in call to Add()")
	}

	if err := verifyURL(item); err != nil {
		return fmt.Errorf("%q is not a valid URL: %w", item, err)
	}

	if slices.Contains(m.Items, item) {
		logger.Warn(fmt.Sprintf("%q already exists in %s", item, m.name))
		return ErrDuplicate
	}

	logger.Info(fmt.Sprintf("adding %q to %s", item, m.name))

	m.Items = append(m.Items, item)
	m.dirty = true

	return nil
}

// Commit writes the list to Firestore
func (m *metaList) Commit() (*gfs.WriteResult, error) {
	if fsClient == nil {
		return nil, errors.New("firestore client not initialized")
	}

	if m.dirty {
		fsCollection, err := getCollection(m.configCollectionString())
		if err != nil {
			return nil, fmt.Errorf("getCollection: %w", err)
		}

		// Split the collection string into collection and doc strings for Firestore
		db := strings.Split(fsCollection, "/")
		collection := db[0]
		doc := strings.Join(db[1:], "/")

		logger.Info(fmt.Sprintf("committing changes to %q", fsCollection))

		return fsClient.Collection(collection).Doc(doc).Set(context.Background(), m)
	}

	return nil, nil
}

// Count returns the number of items in the list
func (m *metaList) Count() int {
	return len(m.Items)
}

// Remove an item from the list
func (m *metaList) Remove(item string) error {
	if item == "" {
		return errors.New("item cannot be blank in call to Remove()")
	}

	logger.Info(fmt.Sprintf("removing %q from %s", item, m.name))

	for i, v := range m.Items {
		if v == item {
			m.Items = slices.Delete(m.Items, i, i)
		}
	}

	m.dirty = true

	// TODO: Remove entries depending on type (modinfo, toolinfo, etc.)
	return nil
}

func (m *metaList) String() string {
	out := []string{}

	for i, item := range m.Items {
		out = append(out, fmt.Sprintf("%04d: %s", i+1, item))
	}
	return strings.Join(out, "\n")
}

func (m *metaList) MarshalJSON() ([]byte, error) {
	out := []byte{}

	if j, err := json.MarshalIndent(m.Items, "  ", "  "); err != nil {
		return out, err
	} else {
		return fmt.Appendf(out, "{\n  %q: %s\n}\n", m.name, string(j)), nil
	}
}

func (m *metaList) configCollectionString() string {
	return metaCollectionBase + "." + m.name
}

/*
// Public Functions
*/
func ModInfo() (DBList[string], error) {
	err := modInfoCache.Fetch()
	if err != nil {
		return nil, fmt.Errorf("Fetch: %w", err)
	}

	return &modInfoCache, nil
}

func Repos() (DBList[string], error) {
	err := repoCache.Fetch()
	if err != nil {
		return nil, fmt.Errorf("Fetch: %w", err)
	}

	return &repoCache, nil
}

func ToolInfo() (DBList[string], error) {
	err := toolInfoCache.Fetch()
	if err != nil {
		return nil, fmt.Errorf("Fetch: %w", err)
	}

	return &toolInfoCache, nil
}
