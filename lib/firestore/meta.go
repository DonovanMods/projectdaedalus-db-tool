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

type MetaList interface {
	fmt.Stringer
	json.Marshaler
	Commit() (*gfs.WriteResult, error)
	Count() int
	Remove(item string) error
	Update(item string, newItem string) error
}

const collectionBase = "firebase.collections.meta"

var repo = metaList{name: "repositories"}
var modInfo = metaList{name: "modinfo"}
var toolInfo = metaList{name: "toolinfo"}

type metaList struct {
	Items []string `firestore:"list"`
	name  string   `firestore:"-" json:"-"`
	dirty bool
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

// Update updates or adds an item to the list
// If item is already in the list, it will be removed and replaced with newItem
// if item is blank, newItem will be added
func (m *metaList) Update(item string, newItem string) error {
	// do not add blank items
	if newItem == "" {
		return errors.New("newItem cannot be blank in call to Update()")
	}

	if item == newItem {
		return fmt.Errorf("%q is the same as %q in %s", item, newItem, m.name)
	}

	if item == "" {
		logger.Info(fmt.Sprintf("adding %q to %s", newItem, m.name))
	} else {
		logger.Info(fmt.Sprintf("updating %q with %q in %s", item, newItem, m.name))

		if slices.Contains(m.Items, item) {
			if err := m.Remove(item); err != nil {
				return err
			}
		}
	}

	if slices.Contains(m.Items, newItem) {
		return fmt.Errorf("%q already exists in %s", newItem, m.name)
	}

	m.Items = append(m.Items, newItem)
	m.dirty = true

	return nil
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

	return nil
}

func (m *metaList) String() string {
	return strings.Join(m.Items, "\n")
}

func (m *metaList) MarshalJSON() ([]byte, error) {
	var out = []byte{}

	if j, err := json.MarshalIndent(m.Items, "  ", "  "); err != nil {
		return out, err
	} else {
		return fmt.Appendf(out, "{\n  %q: %s\n}\n", m.name, string(j)), nil
	}
}

func (m *metaList) configCollectionString() string {
	return collectionBase + "." + m.name
}

func ModInfo() (MetaList, error) {
	if modInfo.Items != nil {
		return &modInfo, nil
	}

	return getDataFor(&modInfo)
}

func Repos() (MetaList, error) {
	if repo.Items != nil {
		return &repo, nil
	}

	return getDataFor(&repo)
}

func ToolInfo() (MetaList, error) {
	if toolInfo.Items != nil {
		return &toolInfo, nil
	}

	return getDataFor(&toolInfo)
}

func getDataFor(structPtr *metaList) (*metaList, error) {
	docSnap, err := getDocument((*structPtr).configCollectionString())
	if err != nil {
		return nil, fmt.Errorf("getDocument: %w", err)
	}

	if !docSnap.Exists() {
		return nil, errors.New("document does not exist")
	}

	if err := docSnap.DataTo(&structPtr); err != nil {
		return nil, fmt.Errorf("DataTo: %w", err)
	}

	return structPtr, nil
}
