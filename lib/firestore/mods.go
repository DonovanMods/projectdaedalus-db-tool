package firestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gfs "cloud.google.com/go/firestore"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/mod"
	"google.golang.org/api/iterator"
)

const modsCollectionBase = "firebase.collections.mods"

type Mods interface {
	fmt.Stringer
	json.Marshaler
	Commit() (*gfs.WriteResult, error)
	Count() int
}

var modList = mods{}

type mods struct {
	Items []mod.Mod
	dirty bool
}

// Commit writes the list to Firestore
func (m *mods) Commit() (*gfs.WriteResult, error) {
	if !m.dirty {
		return nil, nil
	}

	if fsClient == nil {
		return nil, errors.New("firestore client not initialized")
	}

	fsCollection, err := getCollection(modsCollectionBase)
	if err != nil {
		return nil, fmt.Errorf("getCollection: %w", err)
	}

	logger.Info(fmt.Sprintf("committing changes to %q", fsCollection))

	doc := fsClient.Collection(fsCollection).NewDoc()

	return doc.Set(context.Background(), m)
}

func (m *mods) Count() int {
	return len(m.Items)
}

func (m *mods) String() string {
	out := []string{}

	for i, item := range m.Items {
		out = append(out, fmt.Sprintf("%04d: %s", i+1, item.String()))
	}
	return strings.Join(out, "\n")
}

func (m *mods) MarshalJSON() ([]byte, error) {
	out := []byte{}

	if j, err := json.MarshalIndent(m.Items, "  ", "  "); err != nil {
		return out, err
	} else {
		return fmt.Appendf(out, "{\n  %q: %s\n}\n", "mods", string(j)), nil
	}
}

/*
// Public Functions
*/
func ModList() (Mods, error) {
	if modList.Items != nil {
		return &modList, nil
	}

	return getMods(&modList)
}

/*
// Private Functions
*/
func getMods(mods *mods) (*mods, error) {
	iter, err := getDocuments(modsCollectionBase)
	if err != nil {
		return nil, fmt.Errorf("getDocument: %w", err)
	}
	defer iter.Stop()

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		fmt.Println(doc.Data())

		m := (mod.Mod{}).New()
		if err := doc.DataTo(&m); err != nil {
			return nil, fmt.Errorf("DataTo: %w", err)
		}
		logger.Info(fmt.Sprintf("retrieved %s", m.Name))
		m.Meta.State = mod.Unmodified
		mods.Items = append(mods.Items, m)
	}

	logger.Info("successfully retrieved mods list")

	return mods, nil
}
