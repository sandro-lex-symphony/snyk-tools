package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"snykctl/snykTool"
	"sort"
	"time"
)

func usage() {
	fmt.Printf("Usage:\n" +
		"\tconfigure\n" +
		"\ttoken\n" +
		"\tlist-users [org] [prj]\n" +
		"\tadd-user [org] [prj]\n" +
		"\tlist-group-users\n" +
		"\tcompare-users [org] [org]\n" +
		"\tcopy-users [org] [org]\n" +
		"\tlist-orgs\n" +
		"\tsearch-org [name]\n" +
		"\tcreate-org [name]\n" +
		"\tlist-projects [org]\n" +
		"\tsearch-projects [org]\n" +
		"\tdelete-project [org] [prj]\n" +
		"\tdelete-all-projects [org]\n" +
		"\tproject [org] [prj]\n" +
		"\tlist-project-issues [org] [prj]\n" +
		"\treport-org-issues [org]\n" +
		"\tlist-group-ignores\n" +
		"\tlist-org-ignores [org]\n" +
		"\tlist-project-ignores [org] [prj]\n" +
		"\tissue-count [org]\n" +
		"\tissue-count [org] [prj]\n" +
		"\tget-org-config [org]\n")
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
	htmlFlag := flag.Bool("html", false, "Html table")
	// TODO: Add ansync request option flag
	flag.Parse()
	if *debugFlag {
		snykTool.SetDebug(true)
	}

	if *htmlFlag {
		snykTool.SetHtmlFormat(true)
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
		fmt.Printf("group_id [.... %s]: ", group_id[len(group_id)-6:])
		text, _ = reader.ReadString('\n')
		if len(text) > 2 {
			group_id = text[:len(text)-1]
		}

		snykTool.WriteConf(token, group_id)
	case "token":
		fmt.Println(snykTool.GetToken())

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

	case "compare-users":
		r1, err := snykTool.ListUsers(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		s1 := snykTool.GetOrgName(os.Args[2])
		r2, err := snykTool.ListUsers(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		s2 := snykTool.GetOrgName(os.Args[3])
		snykTool.FormatUsers2Cols(r1, r2, s1, s2)

	case "copy-users":
		snykTool.CopyUsers(os.Args[2], os.Args[3])
		fmt.Println("OK")

	case "add-user":
		snykTool.AddUser(os.Args[2], os.Args[3], "collaborator")
		fmt.Println("OK")

	case "get-org-config":
		snykTool.GetOrgConfig(flag.Arg(1))
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

	case "project":
		snykTool.GetProject(flag.Arg(1), flag.Arg(2))

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

	case "delete-project":
		if snykTool.DeleteProject(flag.Arg(1), flag.Arg(2)) {
			fmt.Println("OK")
		}

	case "delete-all-projects":
		prjs, err := snykTool.GetProjects(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		for _, prj := range prjs.Projects {
			if snykTool.DeleteProject(flag.Arg(1), prj.Id) {
				fmt.Printf("%s DELETED\n", prj.Id)
			}
		}
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
			if !snykTool.IsQuiet() {
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
		for _, issue := range result.Issues {
			if *quietFlag {
				fmt.Printf("%s\n", issue.Id)
			} else {
				fmt.Printf("%s\t%s\t%s\n", issue.Id, issue.PkgName, issue.IssueData.Severity)
			}
		}
	case "compare-project-issues":
		result1, err := snykTool.GetProjectIssues(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}
		result2, err := snykTool.GetProjectIssues(flag.Arg(3), flag.Arg(4))
		if err != nil {
			log.Fatal(err)
		}

		var res1, res2 []string
		for _, issue := range result1.Issues {
			res1 = append(res1, issue.Id)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(res1)))
		for _, issue := range result2.Issues {
			res2 = append(res2, issue.Id)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(res2)))

		s3 := merge(res1, res2)
		for i := 0; i < len(s3); i++ {
			if contains(res1, s3[i]) && contains(res2, s3[i]) {
				fmt.Printf("%s\t\t\t\t%s\n", s3[i], s3[i])
			} else if contains(res1, s3[i]) {
				fmt.Printf("%s\t\t\t\t\t------MISSING------\n", s3[i])
			} else if contains(res2, s3[i]) {
				fmt.Printf("------MISSING------\t\t\t\t\t%s\n", s3[i])
			} else {
				fmt.Printf("===== ERROR ====\n")
			}
		}

	case "issue-count":
		if flag.Arg(2) == "" {
			org_issue_count(flag.Arg(1))
		} else {
			prj_issue_count(flag.Arg(1), flag.Arg(2))
		}

	case "report-org-issues":
		// get all prjs for the  org
		// for each prj get all the issues
		// count sevs
		result, err := snykTool.GetProjects(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		var c int
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
			for _, issue := range res.Issues {
				if "critical" == issue.IssueData.Severity {
					c += 1
				} else if "high" == issue.IssueData.Severity {
					h += 1
				} else if "medium" == issue.IssueData.Severity {
					m += 1
				} else {
					l += 1
				}
			}
		}
		fmt.Println("P:", prjs)
		fmt.Println("C: ", c)
		fmt.Println("H: ", h)
		fmt.Println("M: ", m)
		fmt.Println("L: ", l)

	}
}

func org_issue_count(org_id string) {
	result, err := snykTool.GetProjects(org_id)
	if err != nil {
		log.Fatal(err)
	}

	var c int
	var h int
	var m int
	var l int
	var prjs int

	htmlTable := "<table border=1><tr><thead><th>Project</th><th>Critical</th><th>High</th><th>Medium</th><th>Low</th></tr></thead><tbody>"

	for _, project := range result.Projects {
		prjs += 1
		r := snykTool.IssuesCount(org_id, project.Id)
		htmlTable += snykTool.FormatPrjIssuesCountHtml(r, org_id, project.Id)
		for _, result := range *r.Results {
			c += result.Severity.Critical
			h += result.Severity.High
			m += result.Severity.Medium
			l += result.Severity.Low
		}

	}

	htmlTable += fmt.Sprintf("<tr><th align=left>TOTAL</th><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", c, h, m, l)
	htmlTable += "</tbody></table>"

	fmt.Print(htmlTable)

}

func prj_issue_count(org_id, prj_id string) {
	result := snykTool.IssuesCount(org_id, prj_id)
	if snykTool.IsHtmlFormat() {
		fmt.Printf("%s", snykTool.FormatPrjIssuesCountHtml(result, flag.Arg(1), flag.Arg(2)))
	} else {
		fmt.Printf("%s", snykTool.FormatIssuesResult(result, flag.Arg(1), flag.Arg(2)))
	}
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false

}

func merge(s1 []string, s2 []string) []string {
	var s3 []string
	s3 = s1
	for i := 0; i < len(s2); i++ {
		if !contains(s1, s2[i]) {
			s3 = append(s3, s2[i])
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(s3)))
	return s3

}
