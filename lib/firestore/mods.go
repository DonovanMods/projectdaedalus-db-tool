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

// Data Cache
var modCache = mods{name: "mods"}

type mods struct {
	Items []mod.Mod
	name  string
	dirty bool
}

func (M *mods) Fetch() error {
	if M.Items != nil {
		logger.Info("Using cached data for mods")
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

	logger.Info("successfully retrieved mods")

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

	fsCollection, err := getCollection(modsCollectionBase)
	if err != nil {
		return nil, fmt.Errorf("getCollection: %w", err)
	}

	logger.Info(fmt.Sprintf("committing changes to %q", fsCollection))

	doc := fsClient.Collection(fsCollection).NewDoc()

	return doc.Set(context.Background(), M)
}

func (M *mods) Count() int {
	return len(M.Items)
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
		return fmt.Appendf(out, "{\n  %q: %s\n}\n", "mods", string(j)), nil
	}
}

// Update updates or adds an item to the list
// If item is already in the list, it will be removed and replaced with newItem
// if item is blank, newItem will be added
func (M *mods) Update(item mod.Mod, newItem mod.Mod) error {
	// FIXME: Implement Udate() for Mods
	logger.Warn("Update() not implemented for mods")
	return nil
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
