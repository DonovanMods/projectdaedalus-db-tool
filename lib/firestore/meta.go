package firestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type MetaList interface {
	Items() []string
	Name() string
	Commit() error
	Remove(item string) error
	String() string
	Json() (string, error)
	Update(item string, newItem string) error
}

var repo = metaList{name: "repositories"}
var modInfo = metaList{name: "modinfo"}
var toolInfo = metaList{name: "toolinfo"}

type metaList struct {
	List []string `firestore:"list"`
	name string
}

// FIXME: Commit() is not implemented
func (m *metaList) Commit() error {
	// return setDocument(fmt.Sprintf("firebase.collections.meta.%s", m.name), m)
	fmt.Println("Commit() not yet implemented -- items will not be written to the database")
	return nil
}

// Update updates or adds an item to the list
// If item is already in the list, it will be removed and replaced with newItem
// if item is blank, newItem will be added
func (m *metaList) Update(item string, newItem string) error {
	// do not add blank items
	if newItem == "" {
		return errors.New("newItem cannot be blank in call to Update()")
	}

	if item != "" && slices.Contains(m.List, item) {
		if err := m.Remove(item); err != nil {
			return err
		}
	}

	m.List = append(m.List, newItem)

	return nil
}

// Remove an item from the list
func (m *metaList) Remove(item string) error {
	if item == "" {
		return errors.New("item cannot be blank in call to Remove()")
	}

	for i, v := range m.List {
		if v == item {
			m.List = slices.Delete(m.List, i, i)
		}
	}

	return nil
}

func (m *metaList) String() string {
	return strings.Join(m.List, "\n")
}

func (m *metaList) Json() (string, error) {
	if j, err := json.MarshalIndent(m.List, "  ", "  "); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("{\n  %q: %s\n}\n", m.name, string(j)), nil
	}
}

// Items returns the slice of items
func (m *metaList) Items() []string {
	return m.List
}

// Name returns the name of the list
func (m *metaList) Name() string {
	return m.name
}

func ModInfo() (MetaList, error) {
	if modInfo.List != nil {
		return &modInfo, nil
	}

	return getDataFor(&modInfo)
}

func Repos() (MetaList, error) {
	if repo.List != nil {
		return &repo, nil
	}

	return getDataFor(&repo)
}

func ToolInfo() (MetaList, error) {
	if toolInfo.List != nil {
		return &toolInfo, nil
	}

	return getDataFor(&toolInfo)
}

func getDataFor(structPtr *metaList) (*metaList, error) {
	docSnap, err := getDocument(fmt.Sprintf("firebase.collections.meta.%s", (*structPtr).name))
	if err != nil {
		return nil, err
	}

	if !docSnap.Exists() {
		return nil, errors.New("document does not exist")
	}

	if err := docSnap.DataTo(&structPtr); err != nil {
		return nil, err
	}

	return structPtr, nil
}
