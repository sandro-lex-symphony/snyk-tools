package snykTool

import (
	"fmt"
)

var Quiet bool
var NameOnly bool

func SetQuiet(b bool){
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
		if IsQuiet(){
			fmt.Printf("%s\n", user.Id)
		} else if IsNameOnly() {
			fmt.Printf("%s\n", user.Email)
		} else {
			fmt.Printf("%s\t%s\t%s\n", user.Id, user.Role, user.Name)
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