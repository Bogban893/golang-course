package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Date        string `json:"created_at"`
}

func request(owner string, repo string) (*RepoInfo, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("repository not found")
		case http.StatusForbidden:
			return nil, fmt.Errorf("access forbidden or rate limit exceeded")
		default:
			return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
		}
	}

	var res RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return &res, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Incorrect \tinput \nExpected \tinput: owner repo")
		return
	}
	res, err := request(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := time.Parse(time.RFC3339, res.Date)
	fmt.Printf("Repository: \t%s\nDescription: \t%s\nStars:  \t%d\nForks:  \t%d\nCreated: \t%s\n",
		res.Name, res.Description, res.Stars, res.Forks, t.Format("January 2, 2006"))
}
