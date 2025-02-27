package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/spf13/viper"
)

type repoList struct {
	List []string `firestore:"list"`
}

func (r *repoList) Add(repo string) {
	r.List = append(r.List, repo)
}

func (r *repoList) Remove(repo string) {
	for i, v := range r.List {
		if v == repo {
			r.List = append(r.List[:i], r.List[i+1:]...)
		}
	}
}

func (r *repoList) Print() {
	for _, v := range r.List {
		fmt.Println(v)
	}
}

func (r *repoList) JSON() string {
	j, _ := json.Marshal(r.List)
	return fmt.Sprintf(`{"repos":%s}`, string(j))
}

var repos *repoList

func Repos() *repoList {
	if repos != nil {
		return repos
	}

	repoCollection := viper.GetString("firebase.collections.meta.repositories")
	if repoCollection == "" {
		log.Fatal("No repository collection specified in config")
	}
	logger.Log.Info(fmt.Sprintf("Fetching repositories from %q", repoCollection))

	client, err := getClient()
	if err != nil {
		log.Fatal(err)
	}

	reposRef := client.Doc(repoCollection)

	docsnap, err := reposRef.Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err := docsnap.DataTo(&repos); err != nil {
		log.Fatal(err)
	}

	return repos
}
