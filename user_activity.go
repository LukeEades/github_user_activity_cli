package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const template_url = "https://api.github.com/users/%s/events"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("syntax: ./user_activity username")
		return
	}
	url := fmt.Sprintf(template_url, os.Args[1])
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("problem with fetching user data")
		return
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	stat_code := response.StatusCode
	if stat_code == 403 || stat_code == 404 {
		fmt.Println("no such user")
		return
	}
	if stat_code == 503 {
		fmt.Println("api unavailable")
		return
	}
	if err != nil {
		fmt.Println("problem with fetching user data")
		return
	}
	var ob interface{}
	err = json.Unmarshal(body, &ob)
	if err != nil {
		fmt.Println("problem unmarshalling json")
		return
	}
	m := ob.([]interface{})
	for _, val := range m {
		n := val.(map[string]any)
		var mess string
		repo_name := n["repo"].(map[string]any)["name"].(string)
		print := true
		switch n["type"] {
		case "WatchEvent":
			{
				mess = "Starred " + n["repo"].(map[string]any)["name"].(string)
			}
		case "PushEvent":
			{
				num := int(n["payload"].(map[string]any)["size"].(float64))
				mess = "Pushed " + fmt.Sprintf("%d", num) + " commit to " + repo_name
				if num > 1 {
					mess = "Pushed " + fmt.Sprintf("%d", num) + " commits to " + repo_name
				}
			}
		case "IssuesEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " an issue in " + repo_name
			}
		case "SponsorshipEvent":
			{
				mess = "Sponsored " + n["org"].(map[string]any)["org"].(map[string]any)["login"].(string)
			}
		case "ReleaseEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " a release of " + repo_name
			}
		case "PullRequestReviewThreadEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " a comment thread on a pull request on " + repo_name
			}
		case "PullRequestReviewCommentEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " a comment on a pull request review on " + repo_name
			}
		case "PullRequestReviewEvent":
			{
				mess = "Reviewed a pull request for " + repo_name
			}
		case "PullRequestEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " a pull request on " + repo_name
			}
		case "PublicEvent":
			{
				mess = "Changed visibility status of " + repo_name + " to public"
			}
		case "MemberEvent":
			{
				mess = "Added " + n["payload"].(map[string]any)["member"].(string) + " to " + repo_name
			}
		case "IssueCommentEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " a comment on an issue in " + repo_name
			}
		case "ForkEvent":
			{
				mess = "Forked " + repo_name
			}
		case "DeleteEvent":
			{
				mess = "Deleted a " + n["payload"].(map[string]any)["ref_type"].(string) + " on " + repo_name
			}
		case "CreateEvent":
			{
				mess = "Created a " + n["payload"].(map[string]any)["ref_type"].(string) + " on " + repo_name
			}
		case "CommitCommentEvent":
			{
				mess = cases.Title(language.English).String(n["payload"].(map[string]any)["action"].(string)) + " a commit comment in " + repo_name
			}
		default:
			{
				print = false
			}
		}
		if print {
			fmt.Println("- " + mess)
		}
	}
}
