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
	case "list-projects":
        result, err := snykTool.GetProjects(flag.Arg(1))
        if err != nil {
            log.Fatal(err)
        }
        for _, prj := range result.Projects {
            if *quietFlag {
                fmt.Printf("%s\n", prj.Id)
            } else {
                fmt.Printf("%s\t%s\n", prj.Id, prj.Name)
            }
        }
	}
}
