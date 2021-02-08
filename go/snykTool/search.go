package snykTool

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
)

var Debug bool
var Timeout int

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

func RequestGet(path string) (*http.Response) {

    timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
    client := http.Client {
        Timeout: timeout,
    }
    req := SnykURL + path
    if IsDebug() {
        fmt.Println(req)
    }
    request, err := http.NewRequest("GET", req, nil)
    token := GetToken()
    request.Header.Set("Authorization", "token " + token)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    }
    return resp
}


func GetGroupMembers() ([]*User, error) {
    group := GetGroupId()
    resp := RequestGet("/group/" + group + "/members")
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
    resp := RequestGet("/org/" + org_id + "/members")

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

func SearchProjects(org_id string, term string) (*ProjectsResult, error) {
    result, err := GetProjects(org_id)
    if err != nil {
        log.Fatal(err)
    }

    var filtered ProjectsResult
    for _, prj := range(result.Projects) {
        if  strings.Contains(strings.ToLower(prj.Name), strings.ToLower(term)) {
           filtered.Projects = append(filtered.Projects, prj)
        }
    }

    return &filtered, nil
}

func GetProjectIgnores(org_id string, prj_id string) ([]IgnoreResult) {
    resp := RequestGet("/org/" + org_id + "/project/" + prj_id + "/ignores")

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
    timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
    client := http.Client {
        Timeout: timeout,
    }

    request, err := http.NewRequest("POST", SnykURL + "/org/" + org_id + "/project/" + prj_id + "/aggregated-issues", nil)
    token := GetToken()
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", "token " + token)
    if err != nil {
        log.Fatal(err)
    }
    resp, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    }
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
    resp := RequestGet("/org/" + org_id + "/projects")

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

func CreateOrg(org_name string) (*CreateOrgResult, error) {
    timeout := time.Duration(time.Duration(GetTimeout()) * time.Second)
    client := http.Client{
        Timeout: timeout,
    }
    jsonValue, _ := json.Marshal(map[string]string{
        "name": org_name,
    })

    request, err := http.NewRequest("POST", SnykURL + "/org", bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    token := GetToken()
    request.Header.Set("Authorization", "token " + token)
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
    return &result, nil
}

func SearchOrgs(term string) (*OrgList, error) {
    result, err := GetOrgs()
    if err != nil {
        log.Fatal(err)
    }
    var filtered OrgList
    for _, org := range(result.Orgs) {
        if  strings.Contains(strings.ToLower(org.Name), strings.ToLower(term)) {
           filtered.Orgs = append(filtered.Orgs, org)
        }
    }
    return &filtered, nil
}

