package snykTool


const SnykURL = "https://snyk.io/api/v1"

type Org struct {
    Id string
    Name string
}

type OrgList struct {
    Orgs []*Org
}

type User struct {
    Id string
    Name string
    Role string
    Username string
}

type CreateOrgResult struct {
    Id string
    name string
    slug string
    url string
    created string
}

type ProjectsResult struct {
    Org
    Projects []*Project
}

type Project struct {
    Name string
    Id string
}

