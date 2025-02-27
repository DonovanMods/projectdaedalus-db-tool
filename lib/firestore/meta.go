package firestore

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MetaList interface {
	Add(item string)
	Remove(item string)
	JSON() string
	Print()
}

var repo = metaList{name: "repositories"}
var modInfo = metaList{name: "modinfos"}
var toolInfo = metaList{name: "toolinfos"}

type metaList struct {
	List []string `firestore:"list"`
	name string
}

func (m *metaList) Add(item string) {
	m.List = append(m.List, item)
}

func (m *metaList) Remove(item string) {
	for i, v := range m.List {
		if v == item {
			m.List = append(m.List[:i], m.List[i+1:]...)
		}
	}
}

func (m *metaList) JSON() string {
	j, _ := json.Marshal(m.List)
	return fmt.Sprintf(`{%q:%s}`, m.name, string(j))
}

func (m *metaList) Print() {
	for _, v := range m.List {
		fmt.Println(v)
	}
}

func ModInfos() (MetaList, error) {
	if modInfo.List != nil {
		return &modInfo, nil
	}

	return getDataFor("firebase.collections.meta.modinfo", &modInfo)
}

func Repos() (MetaList, error) {
	if repo.List != nil {
		return &repo, nil
	}

	return getDataFor("firebase.collections.meta.repositories", &repo)
}

func ToolInfos() (MetaList, error) {
	if toolInfo.List != nil {
		return &toolInfo, nil
	}

	return getDataFor("firebase.collections.meta.toolinfo", &toolInfo)
}

func getDataFor(collection string, structPtr *metaList) (*metaList, error) {
	docSnap, err := getDocument(collection)
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
