package mod

import (
	"fmt"
)

type mState int

var StateString = map[mState]string{
	StateUnmodified: "Unmodified",
	StateNew:        "New",
	StateUpdated:    "Updated",
	StateDeleted:    "Deleted",
}

const (
	StateUnmodified mState = iota
	StateNew
	StateUpdated
	StateDeleted
)

/*
	 modinfo format:
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
	ID            string `firestore:"-" json:"id"`
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
		Status struct {
			Errors   []string `firestore:"errors" json:"errors"`
			Warnings []string `firestore:"warnings" json:"warnings"`
		} `firestore:"status" json:"status"`
		state mState `firestore:"-" json:"-"`
	} `firestore:"meta" json:"meta"`
}

func (m Mod) New() Mod {
	mod := Mod{}
	mod.Meta.Status.Errors = []string{}
	mod.Meta.Status.Warnings = []string{}
	mod.SetState(StateNew)

	return mod
}

func (m *Mod) String() string {
	return fmt.Sprintf("%s v%s by %s", m.Name, m.Version, m.Author)
}

func (m *Mod) State() mState {
	return m.Meta.state
}

func (m *Mod) SetState(state mState) {
	m.Meta.state = state
}

func (m *Mod) StateString() string {
	return StateString[m.Meta.state]
}

func (m *Mod) Warnings() []string {
	return m.Meta.Status.Warnings
}

func (m *Mod) AddWarning(s string) {
	m.Meta.Status.Warnings = append(m.Meta.Status.Warnings, s)
}

func (m *Mod) Errors() []string {
	return m.Meta.Status.Errors
}

func (m *Mod) AddError(s string) {
	m.Meta.Status.Errors = append(m.Meta.Status.Errors, s)
}

func (m *Mod) Dirty() bool {
	return m.Meta.state != StateUnmodified
}

func (m *Mod) Valid() bool {
	return len(m.Meta.Status.Errors) == 0
}
