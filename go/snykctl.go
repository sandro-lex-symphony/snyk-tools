package main

import (
    "bufio"
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
	switch os.Args[1] {
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
            fmt.Printf("%s\t%s\n", org.Id, org.Name)
        }
	case "search-org":
		result, err := snykTool.SearchOrgs(os.Args[2])
		if err != nil {
            log.Fatal(err)
        }
		for _, org := range result.Orgs {
			fmt.Printf("%s\t%s\n", org.Id, org.Name)
		}
	case "list-projects":
        result, err := snykTool.GetProjects(os.Args[2])
        if err != nil {
            log.Fatal(err)
        }
        for _, prj := range result.Projects {
            fmt.Printf("%s\t%s\n", prj.Id, prj.Name)
        }
	}
}
