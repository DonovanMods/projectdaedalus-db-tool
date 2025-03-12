package mod

import (
	"encoding/json"
	"fmt"
)

type (
	mState  int
	mStatus map[string][]string
)

const (
	Unmodified mState = iota
	New
	Updated
	Deleted
)

/*
	{
	  "name": "First Mod Name",
	  "author": "whatever name you want as the Author",
	  "version": "1.0",
	  "compatibility": "w57",
	  "description": "A description of what your mod does",
	  "files": {
	    "pak": "the direct download URL for your PAK mod file",
	    "exmodz": "the direct download URL for your EXMODZ mod file"
	  },
	  "imageURL": "A direct download URL to an image that will be displayed in the mod list (optional)",
	  "readmeURL": "A link to the 'raw' version of your mod's README"
	},
*/
type Mod struct {
	docID         string `firestore:"-" json:"-"`
	Name          string `firestore:"name" json:"name"`
	Author        string `firestore:"author" json:"author"`
	Version       string `firestore:"version" json:"version"`
	Compatibility string `firestore:"compatibility" json:"compatibility"`
	Description   string `firestore:"description" json:"description"`
	ImageURL      string `firestore:"imageURL" json:"imageURL" omitEmpty:"true"`
	ReadmeURL     string `firestore:"readmeURL" json:"readmeURL" omitEmpty:"true"`
	Files         struct {
		Pak    string `firestore:"pak" json:"pak" omitEmpty:"true"`
		Exmodz string `firestore:"exmodz" json:"exmodz" omitEmpty:"true"`
	} `firestore:"files" json:"files"`
	Meta struct {
		Status mStatus `firestore:"status" json:"status"`
		State  mState  `firestore:"state" json:"state"`
	} `firestore:"meta" json:"meta"`
}

func (m Mod) New() Mod {
	mod := Mod{}
	// Our firestore DB assumes these two arrays are always present
	mod.Meta.Status = mStatus{
		"Errors":   []string{},
		"Warnings": []string{},
	}
	mod.Meta.State = New

	return mod
}

func (m *Mod) ID() string {
	return m.docID
}

func (m *Mod) String() string {
	return fmt.Sprintf("%s v%s by %s", m.Name, m.Version, m.Author)
}

func (m *Mod) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Mod) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Mod) State() string {
	switch m.Meta.State {
	case Unmodified:
		return "Unmodified"
	case New:
		return "New"
	case Updated:
		return "Updated"
	case Deleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

func (m *Mod) Dirty() bool {
	return m.Meta.State != Unmodified
}
