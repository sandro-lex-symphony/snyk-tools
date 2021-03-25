package snykTool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var Debug bool
var Timeout int
var OrgsCache *OrgList

func SetTimeout(t int) {
	Timeout = t
}

func GetTimeout() int {
	if Timeout > 0 {
		return Timeout
	}
	return 10
}

func SetDebug(d bool) {
	Debug = d
}

func IsDebug() bool {
	return Debug
}

func RequestGet(path string) *http.Response {
	timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req := SnykURL + path
	if IsDebug() {
		fmt.Println("GET", req)
	}
	request, err := http.NewRequest("GET", req, nil)
	request.Header.Set("Authorization", "token "+GetToken())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func RequestPost(path string, data []byte) *http.Response {
	timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req := SnykURL + path
	if IsDebug() {
		fmt.Println("POST", req)
	}

	request, err := http.NewRequest("POST", req, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+GetToken())

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func GetGroupMembers() ([]*User, error) {
	path := fmt.Sprintf("/group/%s/members", GetGroupId())
	resp := RequestGet(path)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GetGroupMembers failed %s", resp.Status)
	}
	var result []*User
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return result, nil
}

func ListUsers(org_id string) ([]*User, error) {
	path := fmt.Sprintf("/org/%s/members", org_id)
	resp := RequestGet(path)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("ListUsers failed %s", resp.Status)
	}

	var result []*User
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return result, nil
}

func AddUser(org_id string, user_id string, role string) {
	path := fmt.Sprintf("/group/%s/org/%s/members", GetGroupId(), org_id)
	jsonValue, _ := json.Marshal(map[string]string{
		"userId": user_id,
		"role":   role,
	})
	resp := RequestPost(path, jsonValue)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatal("Add User failed ", resp.Status)
	}
}

func CopyUsers(o1, o2 string) {
	// copy all users from o1 to o2
	// dont check if it exists or if already present
	// get users o1
	// add to o2
	result, err := ListUsers(o1)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range result {
		AddUser(o2, user.Id, "collaborator")
	}
}

func SearchProjects(org_id string, term string) (*ProjectsResult, error) {
	result, err := GetProjects(org_id)
	if err != nil {
		log.Fatal(err)
	}

	var filtered ProjectsResult
	for _, prj := range result.Projects {
		if strings.Contains(strings.ToLower(prj.Name), strings.ToLower(term)) {
			filtered.Projects = append(filtered.Projects, prj)
		}
	}

	return &filtered, nil
}

func GetProjectIgnores(org_id string, prj_id string) []IgnoreResult {
	path := fmt.Sprintf("/org/%s/project/%s/ignores", org_id, prj_id)
	resp := RequestGet(path)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatal("Get Ignores failed ", resp.Status)
	}

	var ignore_result map[string][]IgnoreStar

	if err := json.NewDecoder(resp.Body).Decode(&ignore_result); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}

	var result []IgnoreResult

	for key, value := range ignore_result {
		for i := 0; i < len(value); i++ {
			var ii IgnoreResult
			ii.Id = key
			ii.Content = value[i].Star
			result = append(result, ii)
		}
	}

	resp.Body.Close()
	return result
}

func GetProjectIssues(org_id string, prj_id string) (*ProjectIssuesResult, error) {
	path := fmt.Sprintf("/org/%s/project/%s/aggregated-issues", org_id, prj_id)
	resp := RequestPost(path, nil)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GetProjectIssues failed %s", resp.Status)
	}

	var result ProjectIssuesResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func GetProjects(org_id string) (*ProjectsResult, error) {
	path := fmt.Sprintf("/org/%s/projects", org_id)
	resp := RequestGet(path)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GetProjects failed %s", resp.Status)
	}

	var result ProjectsResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return &result, nil
}

func GetProject(org_id, prj_id string) error {
	path := fmt.Sprintf("/org/%s/project/%s", org_id, prj_id)
	resp := RequestGet(path)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("GetProjects failed %s", resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(bytes))
	return nil
}

func CreateOrg(org_name string) (*CreateOrgResult, error) {
	timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	jsonValue, _ := json.Marshal(map[string]string{
		"name": org_name,
	})

	request, err := http.NewRequest("POST", SnykURL+"/org", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	token := GetToken()
	request.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return nil, fmt.Errorf("CreateOrg failed %s", resp.Status)
	}

	var result CreateOrgResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil

}

func GetOrgs() (*OrgList, error) {
	if OrgsCache != nil {
		return OrgsCache, nil
	}

	resp := RequestGet("/orgs")

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GetOrgs failed %s", resp.Status)
	}

	var result OrgList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	OrgsCache = &result
	return &result, nil
}

func GetOrgName(id string) string {
	orgs, err := GetOrgs()
	if err != nil {
		log.Fatal(err)
	}
	for _, o := range orgs.Orgs {
		if o.Id == id {
			return o.Name
		}
	}
	return ""
}

func SearchOrgs(term string) (*OrgList, error) {
	result, err := GetOrgs()
	if err != nil {
		log.Fatal(err)
	}
	var filtered OrgList
	for _, org := range result.Orgs {
		if strings.Contains(strings.ToLower(org.Name), strings.ToLower(term)) {
			filtered.Orgs = append(filtered.Orgs, org)
		}
	}
	return &filtered, nil
}

func IssuesCount(org_id, prj_id string) IssuesResults {
	path := "/reporting/counts/issues/latest?groupBy=severity"
	var str string
	if prj_id == "" {
		str = fmt.Sprintf("{\"filters\": { \"orgs\": [\"%s\"]}}", org_id)
	} else {
		str = fmt.Sprintf("{\"filters\": { \"orgs\": [\"%s\"], \"projects\": [\"%s\"]}}", org_id, prj_id)
	}
	var jsonStr = []byte(str)
	resp := RequestPost(path, jsonStr)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatal("Issue count failed ", resp.Status)
	}
	var result IssuesResults
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	return result
}
