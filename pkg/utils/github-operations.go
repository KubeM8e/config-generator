package utils

import (
	"bytes"
	"config-generator/models"
	"encoding/json"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO: get this from env
const (
	accessToken = "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"
	githubAPI   = "https://api.github.com"
)

func CreateGitHubRepo(repoName string) {

	repo := models.GitHubRepoData{
		Name:        repoName,
		Description: "This is an auto generated repo",
		Private:     false,
		AutoInit:    true,
	}

	b, err := json.Marshal(repo)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", githubAPI+"/user/repos", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "token "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	//if res.StatusCode != 201 {
	//	fmt.Println("Error creating repository:", res.Status)
	//	return
	//}

	log.Println("Repository created successfully")

}

func CloneGitHubRepo(repoName string, tempFolder string) (*git.Worktree, *git.Repository) {
	url := "https://github.com/Shenali-SJ/" + repoName + ".git"

	// TODO: replace with organization details - the username
	repository, err := git.PlainClone(tempFolder, false, &git.CloneOptions{
		URL: url,
		Auth: &gitHttp.BasicAuth{
			Username: "Shenali-SJ",
			Password: accessToken,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Could not cloan - %v", err)
	}

	workTree, err := repository.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	return workTree, repository
}

func PushToGitHub(workTree *git.Worktree, repo *git.Repository, files []string) {
	//status, _ := w.Status()

	//person := Person{
	//	Name: "Shenali",
	//	Age:  22,
	//}
	//file, err := os.Create("tmp/application.yaml")
	//if err != nil {
	//	log.Fatal("Error: ", err)
	//}
	//out, _ := yaml.Marshal(&person)
	//file.Write(out)

	for _, file := range files {
		_, err := workTree.Add(file)
		if err != nil {
			log.Printf("Could not add the file: %v\n", err)
		}
	}

	// Commit the changes
	// TODO: replace with organization details
	_, err := workTree.Commit("Add file", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Shenali",
			Email: "shenalijayakody@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Printf("Could not commit: %s", err)
	}

	// Push the changes
	// TODO: replace with organization details - the username
	errPush := repo.Push(&git.PushOptions{
		Auth: &gitHttp.BasicAuth{
			Username: "Shenali-SJ",
			Password: accessToken,
		},
	})
	if errPush != nil {
		log.Printf("Could not push: %s", errPush)
	}

	// TODO: remove?
	//log.Print("Changes pushed successfully")
}
