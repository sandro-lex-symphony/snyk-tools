package snykTool

import (
	"fmt"
	"sort"
)

var Quiet bool
var NameOnly bool

func SetQuiet(b bool) {
	Quiet = b
}

func IsQuiet() bool {
	return Quiet
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

func FormatIssuesResult(r IssuesResults, org_id, prj_id string) {
	for _, result := range *r.Results {
		fmt.Printf("Org: %s\n", GetOrgName(org_id))
		if prj_id != "" {
			fmt.Printf("Prj: %s\n", prj_id)
		}
		fmt.Printf("Total: %d\nHigh: %d\nMedium: %d\nLow: %d\n", result.Count, result.Severity.High, result.Severity.Medium, result.Severity.Low)
	}
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
