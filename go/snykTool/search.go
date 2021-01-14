package snykTool

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
)

func ListUsers(org_id string) ([]*User, error) {
    timeout := time.Duration(5 * time.Second)
    client := http.Client {
        Timeout: timeout,
    }
    request, err := http.NewRequest("GET", SnykURL + "/org/" + org_id + "/members", nil)
    token := GetToken()
    request.Header.Set("Authorization", "token " + token)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Do(request)
    if err != nil {
        return nil, err
    }

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

func GetProjects(org_id string) (*ProjectsResult, error) {
    timeout := time.Duration(5 * time.Second)
    client := http.Client {
        Timeout: timeout,
    }
    request, err := http.NewRequest("GET", SnykURL + "/org/" + org_id + "/projects", nil)
    token := GetToken()
    request.Header.Set("Authorization", "token " + token)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Do(request)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("ListUsers failed %s", resp.Status)
    }

    var result ProjectsResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        resp.Body.Close()
        return nil, err
    }

    resp.Body.Close()
    return &result, nil
}


func GetOrgs() (*OrgList, error) {
    timeout := time.Duration(5 * time.Second)
    client := http.Client{
        Timeout: timeout,
    }
    request, err := http.NewRequest("GET", SnykURL + "/orgs", nil)
    token := GetToken()
    request.Header.Set("Authorization", "token " + token)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Do(request)
    if err != nil {
        return nil, err
    }

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

