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
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
)

const modsCollectionBase = "firebase.collections.mods"

// Data Cache
var modCache = mods{name: viper.GetString(modsCollectionBase)}

type mods struct {
	Items []mod.Mod
	name  string
	dirty bool
}

func (M *mods) Fetch() error {
	if M.Items != nil {
		logger.Info(fmt.Sprintf("Using cached data for %s", M.name))
		return nil
	}

	iter, err := getDocuments(modsCollectionBase)
	if err != nil {
		return fmt.Errorf("getDocument: %w", err)
	}
	defer iter.Stop()

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}

		m := (mod.Mod{}).New()
		if err := doc.DataTo(&m); err != nil {
			return fmt.Errorf("DataTo: %w", err)
		}
		logger.Info(fmt.Sprintf("retrieved %s", m.Name))
		m.SetState(mod.StateUnmodified)
		m.ID = doc.Ref.ID
		M.Items = append(M.Items, m)
	}

	logger.Info(fmt.Sprintf("successfully retrieved %s", M.name))

	return nil
}

func (M *mods) Add(item mod.Mod) error {
	if item.Name == "" {
		return errors.New("item cannot be blank in call to Add()")
	}

	logger.Info(fmt.Sprintf("adding %q to mods", item.Name))

	// FIXME: Implement Add() for Mods
	logger.Warn("Add() not implemented for mods")

	return nil
}

// Commit writes the list to Firestore
func (M *mods) Commit() (*gfs.WriteResult, error) {
	if !M.dirty {
		return nil, nil
	}

	if fsClient == nil {
		return nil, errors.New("firestore client not initialized")
	}

	// Retrieve the collection name from our config
	collectionName, err := getCollection(modsCollectionBase)
	if err != nil {
		return nil, fmt.Errorf("getCollection: %w", err)
	}

	collection := fsClient.Collection(collectionName)
	if collection == nil {
		return nil, fmt.Errorf("%s collection not found", M.name)
	}

	logger.Info(fmt.Sprintf("committing changes to %q", collectionName))

	for _, m := range M.Items {
		switch m.State() {
		case mod.StateNew:
			if _, err := collection.NewDoc().Set(context.Background(), m); err != nil {
				return nil, err
			}
		case mod.StateUpdated:
			if _, err := collection.Doc(m.ID).Set(context.Background(), m); err != nil {
				return nil, err
			}
		case mod.StateDeleted:
			if _, err := collection.Doc(m.ID).Delete(context.Background()); err != nil {
				return nil, err
			}
		default: // Unmodified
			continue
		}
	}

	return nil, nil
}

func (M *mods) Count() int {
	return len(M.Items)
}

func (M *mods) Parse(item string) error {
	// FIXME: Implement Parse() for Mods
	logger.Warn("Parse() not implemented for mods")
	return nil
}

// Remove an item from the list
func (M *mods) Remove(item mod.Mod) error {
	// FIXME: Implement Remove() for Mods
	logger.Warn("Remove() not implemented for mods")
	return nil
}

func (M *mods) String() string {
	out := []string{}

	for i, item := range M.Items {
		out = append(out, fmt.Sprintf("%04d: %s", i+1, item.String()))
	}
	return strings.Join(out, "\n")
}

func (M *mods) MarshalJSON() ([]byte, error) {
	out := []byte{}

	if j, err := json.MarshalIndent(M.Items, "  ", "  "); err != nil {
		return out, err
	} else {
		return fmt.Appendf(out, "{\n  %q: %s\n}\n", M.name, string(j)), nil
	}
}

/*
// Public Functions
*/
func ModList() (DBList[mod.Mod], error) {
	err := modCache.Fetch()
	if err != nil {
		return nil, fmt.Errorf("Fetch: %w", err)
	}

	return &modCache, nil
}
