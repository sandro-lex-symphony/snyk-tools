package snykctl

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
var ParallelHttpRequests bool
var OrgsCache *OrgList
var FilterLifecycle string
var FilterEnvironment string
var WorkerSize int

func SetWorkerSize(w int) {
	WorkerSize = w
}

func GetWorkerSize() int {
	if WorkerSize > 0 {
		return WorkerSize
	}
	return GetWorkerSizeFromConf()
}

func SetTimeout(t int) {
	Timeout = t
}

// in case the value was passed in CLI options, ignore the conf
func GetTimeout() int {
	if Timeout > 0 {
		return Timeout
	}
	return GetTimeoutFromConf()
}

func SetDebug(d bool) {
	Debug = d
}

func IsDebug() bool {
	return Debug
}

func SetParallelHttpRequests(b bool) {
	ParallelHttpRequests = b
}

func SetFilterLifecycle(lf string) {
	if lf == "dev" {
		lf = "development"
	}
	if lf == "prod" {
		lf = "production"
	}

	FilterLifecycle = lf
}

func SetFilterEnvironment(env string) {
	if env == "front" {
		env = "frontend"
	}
	if env == "back" {
		env = "backend"
	}

	FilterEnvironment = env
}

func MakeGetRequest(path string, ch chan<- *http.Response) {

	// start := time.Now()
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

	// secs := time.Since(start).Seconds()

	ch <- resp
}

func Request(path string, verb string) *http.Response {
	timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req := SnykURL + path
	if IsDebug() {
		fmt.Println(verb, req)
	}
	request, err := http.NewRequest(verb, req, nil)
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

func RequestGet(path string) *http.Response {
	return Request(path, "GET")
}

func RequestDelete(path string) *http.Response {
	return Request(path, "DELETE")
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

func DeleteProject(org_id string, prj_id string) bool {
	path := fmt.Sprintf("/org/%s/project/%s", org_id, prj_id)
	resp := RequestDelete(path)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatal("Get Ignores failed ", resp.Status)
	}
	return true
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

// same as get project, but apply filters
func GetFilteredProjects(org_id string) (*ProjectsResult, error) {
	path := fmt.Sprintf("/org/%s/projects", org_id)

	var filter string
	filter = "{\"filters\": { \"attributes\": {"

	if FilterLifecycle != "" {
		// filter = fmt.Sprintf("{\"filters\": { \"attributes\": { \"lifecycle\": [ \"%s\" ] } } }", FilterLifecycle)
		filter += fmt.Sprintf("\"lifecycle\": [ \"%s\" ]", FilterLifecycle)
	}
	if FilterLifecycle != "" && FilterEnvironment != "" {
		filter += ","
	}

	if FilterEnvironment != "" {
		filter += fmt.Sprintf("\"environment\": [ \"%s\" ]", FilterEnvironment)
	}

	filter += "} } }"

	var jsonStr = []byte(filter)
	resp := RequestPost(path, jsonStr)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatal("Get filtered projects list failed ", resp.Status)
	}
	var result ProjectsResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}

	resp.Body.Close()
	return &result, nil
}

func GetProjects(org_id string) (*ProjectsResult, error) {
	// check if any filter set
	// TODO: generic filters
	if FilterLifecycle != "" || FilterEnvironment != "" {
		return GetFilteredProjects(org_id)
	}

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

func GetOrgConfig(org_id string) {
	path := fmt.Sprintf("/org/%s/settings", org_id)
	resp := RequestGet(path)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		fmt.Printf("GetProjects failed %s", resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(bytes))

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

func GetPrjName(org_id, prj_id string) string {
	path := fmt.Sprintf("/org/%s/project/%s", org_id, prj_id)
	resp := RequestGet(path)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "ERROR"
	}

	var result Project
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return "ERROR"
	}
	resp.Body.Close()

	return result.Name
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

func MakeIssuesCount(org_id, prj_id string, ch chan<- IssuesResults) {
	path := "/reporting/counts/issues/latest?groupBy=severity"
	var str string

	str = fmt.Sprintf("{\"filters\": "+
		"{ \"orgs\": [\"%s\"], "+
		"\"projects\": [\"%s\"], "+
		// "\"ignored\": true, "+
		"\"severity\": [\"critical\",\"high\",\"medium\",\"low\"] "+
		"}}", org_id, prj_id)

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

	ch <- result
}

func IssuesCount(org_id, prj_id string) IssuesResults {
	path := "/reporting/counts/issues/latest?groupBy=severity"
	var str string
	if prj_id == "" {
		str = fmt.Sprintf("{\"filters\": { \"orgs\": [\"%s\"], \"severity\": [\"critical\",\"high\",\"medium\",\"low\"]}}", org_id)
	} else {
		str = fmt.Sprintf("{\"filters\": { \"orgs\": [\"%s\"], \"projects\": [\"%s\"], \"severity\": [\"critical\",\"high\",\"medium\",\"low\"]}}", org_id, prj_id)
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

// could have passed direct by IssueCount with only a org filter
// but then it is not possible to apply other filters like attributes
// so it is better to manually get the list of pojects, consider the filter
// and then get the issue count for each project
// todo: add asnyc parallel requests
func OrgIssueCount(org_id string) []AggregateIssuesResult {
	// this entry should consider filters
	result, err := GetProjects(org_id)
	if err != nil {
		log.Fatal(err)
	}

	var prjs []AggregateIssuesResult

	if ParallelHttpRequests {

		var full_list []string
		for _, v := range result.Projects {
			full_list = append(full_list, v.Id)
		}

		chunks := sliceChunks(full_list, GetWorkerSize())

		ch := make(chan IssuesResults)

		for _, chunk := range chunks {
			for _, prj_id := range chunk {
				go MakeIssuesCount(org_id, prj_id, ch)
			}

			for _, prj_id := range chunk {
				prj := AggregateIssuesResult{
					IssuesResults: <-ch,
					Org:           org_id,
					Prj:           prj_id,
				}
				prjs = append(prjs, prj)
			}
			// time.Sleep(1 * time.Second)
		}

	} else {

		for _, project := range result.Projects {
			prj := AggregateIssuesResult{
				IssuesResults: IssuesCount(org_id, project.Id),
				Org:           org_id,
				Prj:           project.Name,
			}

			prjs = append(prjs, prj)
		}
	}

	return prjs

}

func sliceChunks(full []string, chunksize int) [][]string {
	var chunks [][]string

	for i := 0; i < len(full); i += chunksize {
		end := i + chunksize

		if end > len(full) {
			end = len(full)
		}

		chunks = append(chunks, full[i:end])
	}

	return chunks
}

// func chunkSlice(slice *ProjectsResult, chunkSize int) [][]string {
// 	var chunks [][]string
// 	var chunk []string

// 	count := 0

// 	for _, v := range slice.Projects {
// 		fmt.Printf("XX %s\n", v.Id)
// 		chunk := append(chunk, v.Id)

// 		if count == chunkSize {
// 			count = 0
// 			chunks = append(chunks, chunk)
// 			chunk = nil
// 		}
// 		count += 1
// 	}

// 	// add the reminders
// 	// chunks = append(chunks, chunk)

// 	return chunks
// 	// for i := 0; i < len(slice); i += chunkSize {
// 	// 	end := i + chunkSize

// 	// 	// necessary check to avoid slicing beyond
// 	// 	// slice capacity
// 	// 	if end > len(slice) {
// 	// 		end = len(slice)
// 	// 	}

// 	chunks = append(chunks, slice[i:end])
// }

// return chunks
// }
