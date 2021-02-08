package main

import (
    "bufio"
    "flag"
	"fmt"
	"log"
	"os"
    "snykTool"
    "time"
)

func usage() {
    fmt.Printf("Usage:\n" +
                "\tconfigure\n" +
                "\tlist-users [org] [prj]\n" +
                "\tlist-group-users\n" +
                "\tlist-orgs\n" +
                "\tsearch-org [name]\n" +
                "\tcreate-org [name]\n" +
                "\tlist-projects [org]\n" +
                "\tsearch-projects [org]\n" +
                "\tlist-project-issues [org] [prj]\n" +
                "\treport-org-issues [org]\n")
}

func main() {
	if len(os.Args) < 2 {
        usage()
		os.Exit(1)
	}
    quietFlag := flag.Bool("q", false, "Quiet output")
    nameOnlyFlag := flag.Bool("n", false, "Names only output")
    debugFlag := flag.Bool("d", false, "Debug http requests")
    timeoutFlag := flag.Int("t", 10, "Http timeout")
    flag.Parse()
    if *debugFlag {
        snykTool.SetDebug(true)
    }

    if *quietFlag {
        snykTool.SetQuiet(true)
    }

    if *nameOnlyFlag {
        snykTool.SetNameOnly(true)
    }

    if *timeoutFlag != 10 {
        snykTool.SetTimeout(*timeoutFlag)
    }

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
        snykTool.FormatUser(result)
        
    case "list-group-users":
        result, err := snykTool.GetGroupMembers()
        if err != nil {
            log.Fatal(err)
        }
        snykTool.FormatUser(result)
  
    case "list-orgs":
	    result, err := snykTool.GetOrgs()
	    if err != nil {
	        log.Fatal(err)
        }
        snykTool.FormatOrg(result)

	case "search-org":
		result, err := snykTool.SearchOrgs(flag.Arg(1))
		if err != nil {
            log.Fatal(err)
        }
        snykTool.FormatOrg(result)

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
        snykTool.FormatProjects(result)
        
    case "search-projects":
        result, err := snykTool.SearchProjects(flag.Arg(1), flag.Arg(2))
        if err != nil {
            log.Fatal(err)
        }
        snykTool.FormatProjects(result)

    case "list-project-ignores":
        res := snykTool.GetProjectIgnores(flag.Arg(1), flag.Arg(2))
        snykTool.FormatProjectIgnore(res)
       
    case "list-org-ignores":
        result, err := snykTool.GetProjects(flag.Arg(1))
        if err != nil {
            log.Fatal(err)
        }
        for _, prj := range result.Projects {
            res := snykTool.GetProjectIgnores(flag.Arg(1), prj.Id)
            if !snykTool.IsQuiet() {
                fmt.Printf("====> %s\n", prj.Name)
            }
            snykTool.FormatProjectIgnore(res)
        }

    case "list-group-ignores":
        // doing sequential because of the rate limiting on snyk api
        result, err := snykTool.GetOrgs()
	    if err != nil {
	        log.Fatal(err)
        }
        for _, org := range result.Orgs {
            if ! snykTool.IsQuiet() {
                fmt.Printf("<<< %s >>>\n", org.Name)
            }
            result_prj, err := snykTool.GetProjects(org.Id)
            if err != nil {
                log.Fatal(err)
            }
            for _, prj := range result_prj.Projects {
                result_ignores := snykTool.GetProjectIgnores(org.Id, prj.Id)
                if !snykTool.IsQuiet() {
                    fmt.Printf("====> %s\n", prj.Name)
                }
                snykTool.FormatProjectIgnore(result_ignores)
            }
            // sleep for rate limit
            time.Sleep(1 * time.Second)
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
