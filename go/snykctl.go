package main

import (
    "bufio"
    "flag"
	"fmt"
	"log"
	"os"
	"snykTool"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing args")
		os.Exit(1)
	}
    quietFlag := flag.Bool("q", false, "Quiet output")
    nameOnlyFlag := flag.Bool("n", false, "Names only output")
    flag.Parse()

	switch flag.Arg(0) {
	case "configure":
        token := snykTool.GetToken()
        group_id := snykTool.GetGroupId()

        reader := bufio.NewReader(os.Stdin)
        fmt.Printf("token [.... %s]: ", token[len(token)-6:])
        text, _ := reader.ReadString('\n')
        if len(text) > 2 {
            token = text[:len(text)-1]
        }
        fmt.Printf("group_id [.... %s]: ", group_id[len(group_id)- 6:])
        text, _ = reader.ReadString('\n')
        if len(text) > 2 {
            group_id = text[:len(text) -1]
        }

        snykTool.WriteConf(token, group_id)
    case "list-users":
        result, err := snykTool.ListUsers(os.Args[2])
        if err != nil {
            log.Fatal(err)
        }
        for _, user := range result {
            fmt.Printf("%s\t%s\t%s\n", user.Id, user.Role, user.Name)
        }
    case "list-group-users":
        result, err := snykTool.GetGroupMembers()
        if err != nil {
            log.Fatal(err)
        }
  
        for _, user := range result {
            if *quietFlag {
                fmt.Printf("%s\n", user.Id)
            } else if *nameOnlyFlag {
                fmt.Printf("%s\n", user.Email)
            } else {  
                fmt.Printf("%s\t%s\n", user.Id, user.Email)
            }
        }
    case "list-orgs":
	    result, err := snykTool.GetOrgs()
	    if err != nil {
	        log.Fatal(err)
        }
        for _, org := range result.Orgs {
            if *quietFlag {
                fmt.Printf("%s\n", org.Id)
            } else {
                fmt.Printf("%s\t%s\n", org.Id, org.Name)
            }
        }
	case "search-org":
		result, err := snykTool.SearchOrgs(flag.Arg(1))
		if err != nil {
            log.Fatal(err)
        }
		for _, org := range result.Orgs {
		    if *quietFlag {
                fmt.Printf("%s\n", org.Id)
            } else {
			    fmt.Printf("%s\t%s\n", org.Id, org.Name)
            }
		}
	case "create-org":
        result, err := snykTool.CreateOrg(flag.Arg(1))
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s\n", result.Id)

	case "list-projects":
        result, err := snykTool.GetProjects(flag.Arg(1))
        if err != nil {
            log.Fatal(err)
        }
        for _, prj := range result.Projects {
            if *quietFlag {
                fmt.Printf("%s\n", prj.Id)
            } else if *nameOnlyFlag {
                fmt.Printf("%s\n", prj.Name)
            } else {
                fmt.Printf("%s\t%s\n", prj.Id, prj.Name)
            }
        }
    case "search-projects":
        result, err := snykTool.SearchProjects(flag.Arg(1), flag.Arg(2))
        if err != nil {
            log.Fatal(err)
        }
        for _, prj := range(result.Projects) {
            if *quietFlag {
                fmt.Printf("%s\n", prj.Id)
            } else {
                fmt.Printf("%s\t%s\n", prj.Id, prj.Name)
            }
        }
    case "list-project-issues":
        result, err := snykTool.GetProjectIssues(flag.Arg(1), flag.Arg(2))
        if err != nil {
            log.Fatal(err)
        }
        for _, issue := range(result.Issues) {
            fmt.Printf("%s\t%s\t%s\n", issue.Id, issue.PkgName, issue.IssueData.Severity)
        }
    case "report-org-issues":
        // get all prjs for the  org
        // for each prj get all the issues
        // count sevs
        result, err := snykTool.GetProjects(flag.Arg(1))
        if err != nil {
            log.Fatal(err)
        }
        var h int
        var m int
        var l int
        var prjs int
        for _, project := range result.Projects {
            prjs += 1
            res, err := snykTool.GetProjectIssues(flag.Arg(1), project.Id)
            if err != nil {
                log.Fatal(err)
            }
            for _, issue := range(res.Issues) {
                if "high" == issue.IssueData.Severity {
                    h += 1
                } else if "medium" == issue.IssueData.Severity {
                    m += 1
                } else {
                    l += 1
                }
            }
        }
        fmt.Println("P:", prjs)
        fmt.Println("H: ", h)
        fmt.Println("M: ", m)
        fmt.Println("L: ", l)

	}
}
