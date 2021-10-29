package snykctl

import (
	"fmt"
	"sort"
)

var Quiet bool
var NameOnly bool
var HtmlFormat bool

const IssuesColSize = 60
const IssuesCol2Size = 10

func SetQuiet(b bool) {
	Quiet = b
}

func IsQuiet() bool {
	return Quiet
}

func SetHtmlFormat(b bool) {
	HtmlFormat = b
}

func IsHtmlFormat() bool {
	return HtmlFormat
}

func SetNameOnly(b bool) {
	NameOnly = b
}

func IsNameOnly() bool {
	return NameOnly
}

func FormatUser(result []*User) {
	for _, user := range result {
		if IsQuiet() {
			fmt.Printf("%s\n", user.Id)
		} else if IsNameOnly() {
			fmt.Printf("%s\n", user.Email)
		} else {
			fmt.Printf("%s\t%s\t%s\n", user.Id, user.Role, user.Name)
		}
	}
}

func FormatIssuesResult(r IssuesResults, org_id, prj_id string) string {
	var out string
	for _, result := range *r.Results {
		out = fmt.Sprintf("Org: %s\n", GetOrgName(org_id))
		if prj_id != "" {
			out += fmt.Sprintf("Prj: %s\n", prj_id)
		}
		out += fmt.Sprintf("Total: %d\nCritical: %d\nHigh: %d\nMedium: %d\nLow: %d\n", result.Count, result.Severity.Critical, result.Severity.High, result.Severity.Medium, result.Severity.Low)
	}

	return out
}

func FormatIssues(issuesList []AggregateIssuesResult) string {
	if HtmlFormat {
		return FormatIssuesHtml(issuesList)
	}

	return FormatIssuesCli(issuesList)
}

func FormatIssuesHtml(issuesList []AggregateIssuesResult) string {
	var out string
	var t, c, h, m, l int = 0, 0, 0, 0, 0

	// header
	out += FormatIssuesResultHeaderHtml()
	// each item
	for _, v := range issuesList {
		for _, result := range *v.IssuesResults.Results {
			t += 1 // useless, could've got the length
			c += result.Severity.Critical
			h += result.Severity.High
			m += result.Severity.Medium
			l += result.Severity.Low
			out += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
				v.Prj, result.Severity.Critical, result.Severity.High, result.Severity.Medium, result.Severity.Low)
		}
	}
	// count
	out += fmt.Sprintf("<tr><td>TOTAL</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", c, h, m, l)

	// footer
	out += "</table>"
	return out
}

func FormatIssuesCli(issuesList []AggregateIssuesResult) string {
	var out string
	var t, c, h, m, l int = 0, 0, 0, 0, 0
	// header
	out += FormatIssuesResultHeaderCli()
	// each item
	for _, v := range issuesList {
		for _, result := range *v.IssuesResults.Results {
			var col, line string
			t += 1 // useless, could've got the length
			c += result.Severity.Critical
			h += result.Severity.High
			m += result.Severity.Medium
			l += result.Severity.Low
			line += fmt.Sprintf("%s", v.Prj)
			line = fillSpaces(line, IssuesColSize, " ")
			col = fmt.Sprintf("%d", result.Severity.Critical)
			col = fillSpaces(col, IssuesCol2Size, " ")
			line += col
			col = fmt.Sprintf("%d", result.Severity.High)
			col = fillSpaces(col, IssuesCol2Size, " ")
			line += col
			col = fmt.Sprintf("%d", result.Severity.Medium)
			col = fillSpaces(col, IssuesCol2Size, " ")
			line += col
			col = fmt.Sprintf("%d", result.Severity.Low)
			col = fillSpaces(col, IssuesCol2Size, " ")
			line += col + "\n"

			// line += fmt.Sprintf("%d\t%d\t%d\t%d\n", result.Severity.Critical, result.Severity.High, result.Severity.Medium, result.Severity.Low)
			out += line
			//out += fmt.Sprintf("%s\t\t\t\t%d\t%d\t%d\t%d\n", v.Prj, result.Severity.Critical, result.Severity.High, result.Severity.Medium, result.Severity.Low)
		}
	}

	// footer
	return out
}

func FormatIssuesResultHeaderCli() string {
	var out, col string
	col = fmt.Sprintf("%s", "PROJECT")
	col = fillSpaces(col, IssuesColSize, " ")
	out += col

	col = fmt.Sprintf("%s", "CRITICAL")
	col = fillSpaces(col, IssuesCol2Size, " ")
	out += col
	col = fmt.Sprintf("%s", "HIGH")
	col = fillSpaces(col, IssuesCol2Size, " ")
	out += col
	col = fmt.Sprintf("%s", "MEDIUM")
	col = fillSpaces(col, IssuesCol2Size, " ")
	out += col
	col = fmt.Sprintf("%s", "LOW")
	col = fillSpaces(col, IssuesCol2Size, " ")
	out += col + "\n"

	// out += "CRITICAL\tHIGH\tMEDIUM\tLOW\n"
	return out
	//return "PROJECT\t\t\t\tCRITICAL\tHIGH\tMEDIUM\tLOW\n"
}

func FormatIssuesResultCli(r IssuesResults, org_id, prj_id string) string {
	var out string
	for _, result := range *r.Results {
		out = fmt.Sprintf("%s\t\t\t\t%d\t%d\t%d\t%d\n", GetPrjName(org_id, prj_id), result.Severity.Critical, result.Severity.High, result.Severity.Medium, result.Severity.Low)
	}
	return out
}

func FormatIssuesResultHeaderHtml() string {
	return "<table><tr><td>PROJECT</td><td>Critical</td><td>High</td><td>Medium</td><td>Low</td></tr>"
}
func FormatPrjIssuesCountHtml(r IssuesResults, org_id, prj_id string) string {
	var line string

	for _, result := range *r.Results {
		line = fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
			GetPrjName(org_id, prj_id), result.Severity.Critical, result.Severity.High, result.Severity.Medium, result.Severity.Low)
	}

	return line
}

func FormatUsers2Cols(r1, r2 []*User, s1, s2 string) {
	colsize := 40
	//mid := addSpaces(1, "<=>")
	mid := ""
	left := fillSpaces(s1, colsize, " ")
	fmt.Printf("%s%s%s\n", left, mid, s2)

	leftBar := fillSpaces("", len(s1), "=")
	leftBar = fillSpaces(leftBar, colsize, " ")
	right := fillSpaces("", len(s2), "=")
	fmt.Printf("%s%s%s\n", leftBar, mid, right)

	r3 := mergeUsers(r1, r2)
	for i := 0; i < len(r3); i++ {
		if containUser(r1, r3[i]) && containUser(r2, r3[i]) {
			fmt.Printf("%s%s%s\n", fillSpaces(r3[i].Name, colsize, " "), mid, r3[i].Name)
		} else if containUser(r1, r3[i]) {
			fmt.Printf("%s%s--- MISSING ---\n", fillSpaces(r3[i].Name, colsize, " "), mid)
		} else {
			fmt.Printf("%s%s%s\n", fillSpaces("--- MISSING ---", colsize, " "), mid, r3[i].Name)
		}
	}
}

func FormatOrg(orgs *OrgList) {
	for _, org := range orgs.Orgs {
		if IsQuiet() {
			fmt.Printf("%s\n", org.Id)
		} else {
			fmt.Printf("%s\t%s\n", org.Id, org.Name)
		}
	}
}

func FormatProjects(prjs *ProjectsResult) {
	for _, prj := range prjs.Projects {
		if IsQuiet() {
			fmt.Printf("%s\n", prj.Id)
		} else if IsNameOnly() {
			fmt.Printf("%s\n", prj.Name)
		} else {
			fmt.Printf("%s\t%s\n", prj.Id, prj.Name)
		}
	}
}

func FormatProjectIgnore(res []IgnoreResult) {
	for i := 0; i < len(res); i++ {
		if IsQuiet() {
			fmt.Printf("%s\n", res[i].Id)
		} else {
			fmt.Printf("%s\t%s\t%s\t%s\t\n", res[i].Id, res[i].Content.Created, res[i].Content.IgnoredBy.Email, res[i].Content.Reason)
		}
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

func containUser(u1 []*User, x *User) bool {
	for _, v := range u1 {
		if v.Id == x.Id {
			return true
		}
	}
	return false
}

func mergeUsers(u1 []*User, u2 []*User) []*User {
	var u3 []*User
	u3 = u1
	for i := 0; i < len(u2); i++ {
		if !containUser(u1, u2[i]) {
			u3 = append(u3, u2[i])
		}

	}
	return u3
}

func addSpaces(size int, filler string) string {
	ret := ""
	for i := 0; i < size; i++ {
		ret += filler
	}
	return ret
}

func fillSpaces(s string, size int, fillerChar string) string {
	if len(s) >= size {
		return s
	}
	filler := addSpaces(size-len(s), fillerChar)
	return s + filler
}
