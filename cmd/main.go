package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/sandro-lex-symphony/snyk-tools/snykctl"
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
		"\tlist-projects [org]   [-lifecycle prod | dev | sandbox ] [-env front | back | mobile | onprem ]\n" +
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
		"\tissue-count [org] [prj] [-lifecycle prod | dev | sandbox ] [-env front | back | mobile | onprem ]\n" +
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
	parallelFlag := flag.Bool("p", false, "(Try to) Run HTTP Requests in parallel")
	workerSizeFlag := flag.Int("w", 10, "Number of HTTP requests per worker")
	htmlFlag := flag.Bool("html", false, "Html table")
	lifecycleFlag := flag.String("lifecycle", "", "prod | dev | sandbox")
	environmentFlag := flag.String("env", "", "front | back | onprem | mobile")

	// TODO: Add ansync request option flag
	flag.Parse()
	if *debugFlag {
		snykctl.SetDebug(true)
	}

	if *htmlFlag {
		snykctl.SetHtmlFormat(true)
	}

	if *quietFlag {
		snykctl.SetQuiet(true)
	}

	if *nameOnlyFlag {
		snykctl.SetNameOnly(true)
	}

	if *timeoutFlag != 10 {
		snykctl.SetTimeout(*timeoutFlag)
	}

	if *parallelFlag {
		snykctl.SetParallelHttpRequests(true)
	}

	if *workerSizeFlag != 10 {
		snykctl.SetWorkerSize(*workerSizeFlag)
	}

	if *lifecycleFlag != "" {
		validOptions := []string{"dev", "development", "prod", "production", "sandbox"}
		if !contains(validOptions, *lifecycleFlag) {
			fmt.Printf("ERROR: invalid lifecycle value\n")
			usage()
			os.Exit(1)
		}

		snykctl.SetFilterLifecycle(*lifecycleFlag)
	}

	if *environmentFlag != "" {
		validOptions := []string{"front", "frontend", "back", "backend", "onprem", "mobile"}
		if !contains(validOptions, *environmentFlag) {
			fmt.Printf("ERROR: Invalid env value\n")
			usage()
			os.Exit(1)
		}

		snykctl.SetFilterEnvironment(*environmentFlag)
	}

	switch flag.Arg(0) {
	case "configure":
		token := snykctl.GetToken()
		group_id := snykctl.GetGroupId()

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

		snykctl.WriteConf(token, group_id)
	case "token":
		fmt.Println(snykctl.GetToken())

	case "list-users":
		result, err := snykctl.ListUsers(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		snykctl.FormatUser(result)

	case "list-group-users":
		result, err := snykctl.GetGroupMembers()
		if err != nil {
			log.Fatal(err)
		}
		snykctl.FormatUser(result)

	case "compare-users":
		r1, err := snykctl.ListUsers(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		s1 := snykctl.GetOrgName(os.Args[2])
		r2, err := snykctl.ListUsers(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		s2 := snykctl.GetOrgName(os.Args[3])
		snykctl.FormatUsers2Cols(r1, r2, s1, s2)

	case "copy-users":
		snykctl.CopyUsers(os.Args[2], os.Args[3])
		fmt.Println("OK")

	case "add-user":
		snykctl.AddUser(os.Args[2], os.Args[3], "collaborator")
		fmt.Println("OK")

	case "get-org-config":
		snykctl.GetOrgConfig(flag.Arg(1))
	case "list-orgs":
		result, err := snykctl.GetOrgs()
		if err != nil {
			log.Fatal(err)
		}
		snykctl.FormatOrg(result)

	case "search-org":
		result, err := snykctl.SearchOrgs(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		snykctl.FormatOrg(result)

	case "create-org":
		result, err := snykctl.CreateOrg(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", result.Id)

	case "project":
		snykctl.GetProject(flag.Arg(1), flag.Arg(2))

	case "list-projects":
		result, err := snykctl.GetProjects(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		snykctl.FormatProjects(result)

	case "search-projects":
		result, err := snykctl.SearchProjects(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}
		snykctl.FormatProjects(result)

	case "delete-project":
		if snykctl.DeleteProject(flag.Arg(1), flag.Arg(2)) {
			fmt.Println("OK")
		}

	case "delete-all-projects":
		prjs, err := snykctl.GetProjects(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		for _, prj := range prjs.Projects {
			if snykctl.DeleteProject(flag.Arg(1), prj.Id) {
				fmt.Printf("%s DELETED\n", prj.Id)
			}
		}
	case "list-project-ignores":
		res := snykctl.GetProjectIgnores(flag.Arg(1), flag.Arg(2))
		snykctl.FormatProjectIgnore(res)

	case "list-org-ignores":
		result, err := snykctl.GetProjects(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		for _, prj := range result.Projects {
			res := snykctl.GetProjectIgnores(flag.Arg(1), prj.Id)
			if !snykctl.IsQuiet() {
				fmt.Printf("====> %s\n", prj.Name)
			}
			snykctl.FormatProjectIgnore(res)
		}

	case "list-group-ignores":
		// doing sequential because of the rate limiting on snyk api
		result, err := snykctl.GetOrgs()
		if err != nil {
			log.Fatal(err)
		}
		for _, org := range result.Orgs {
			if !snykctl.IsQuiet() {
				fmt.Printf("<<< %s >>>\n", org.Name)
			}
			result_prj, err := snykctl.GetProjects(org.Id)
			if err != nil {
				log.Fatal(err)
			}
			for _, prj := range result_prj.Projects {
				result_ignores := snykctl.GetProjectIgnores(org.Id, prj.Id)
				if !snykctl.IsQuiet() {
					fmt.Printf("====> %s\n", prj.Name)
				}
				snykctl.FormatProjectIgnore(result_ignores)
			}
			// sleep for rate limit
			time.Sleep(1 * time.Second)
		}

	case "list-project-issues":
		result, err := snykctl.GetProjectIssues(flag.Arg(1), flag.Arg(2))
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
		result1, err := snykctl.GetProjectIssues(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}
		result2, err := snykctl.GetProjectIssues(flag.Arg(3), flag.Arg(4))
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
		result, err := snykctl.GetProjects(flag.Arg(1))
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
			res, err := snykctl.GetProjectIssues(flag.Arg(1), project.Id)
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
	// get the list of issues per org
	// print output according to format options
	aggregatedIssues := snykctl.OrgIssueCount(org_id)
	out := snykctl.FormatIssues(aggregatedIssues)
	fmt.Print(out)
}

func prj_issue_count(org_id, prj_id string) {
	result := snykctl.IssuesCount(org_id, prj_id)
	if snykctl.IsHtmlFormat() {
		fmt.Printf("%s", snykctl.FormatPrjIssuesCountHtml(result, flag.Arg(1), flag.Arg(2)))
	} else {
		fmt.Printf("%s", snykctl.FormatIssuesResultHeaderCli())
		fmt.Printf("%s", snykctl.FormatIssuesResultCli(result, flag.Arg(1), flag.Arg(2)))
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
