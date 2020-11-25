# snyk-tools

## snyk-users
USAGE:         
    snyk-tool list-users org_id [-f] all|json				List users from org_id         
    snyk-tool copy-users src dst					Copy Users from org src to org dst         
    snyk-tool compare-users [src] [dst] -f all				Compare user list from org src and dst         
    snyk-tool search-org [name] 					Search org by name         
    snyk-tool create-org [org_name] 					Create a new org with name [org_name]        
    snyk-tool search-prj [org_id] [-o] origin [-n] name [-f] fmt	Search a project in the org        
    snyk-tool delete-prj [org_id] [prj_id]				Delete project [prj_id] in the org [org_id]]        
    snyk-tool prj-issues [org_id] [prj_id]				Get aggregated issue for [prj_id] in the org [org_id]]        
    snyk-tool configure
