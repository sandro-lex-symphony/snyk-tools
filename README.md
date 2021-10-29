# snykctl

**A cmd line tool for interacting with snyk API**

get the list of orgs:
```bash
snykctl list-orgs
6c5fdc1f-e675-4321-9bb0-bd7a22a34a52    Org-1
f6910fd7-43a3-4e20-8327-6b621b7746b4    Org-2
ecd201fd-2bf1-4ef0-b4b6-2989010b5d48    Org-3 
...
```

search for a specific org
```bash
snykctl search-org rele
8b587d86-66b7-4947-b98c-0242de8b70ce    Release-1
07dfd36f-9193-40c2-8e44-2d5970231bab    release 2020-10-09
de6ec69b-5528-4e2b-bf20-14f6293cb274    relegation 
...
```

list projects for a project
```bash
snykctl list-projects 8b587d86-66b7-4947-b98c-0242de8b70qc
3801e440-e69c-4387-a3a3-9c0a4f2f69fa    com.example.backend.core:malware-scan-client
8c935d18-5de4-4cb7-90fb-4de229f73be6    com.example.backend.core:cache
3458efe9-7c3a-45fd-9138-d62661e572bf    com.example.backend.core:commons
...
```

## Basic Features
**Manipulate API resources**
* list, search create and delete operations works on orgs and projects. 
* show projects informations
* show project config

**Manipulate users**
* list users from a project
* add users to a project
* compare users from two projects
* copy users from one project to another

**Issues**
* shows projects issues list
* Issue count
* Issue report

**Ignores**
* list project ignores
* list org ignores


## Instalation
**Requirements**
* golang > 1.17
* Makefile

```bash
make install
```


## Configure
snykctl users a configuration file located on ~/.snykctl.conf. 
```bash
[DEFAULT]
token = <SNYK_API_TOKEN>
id = <ORG_ID>|<GROUP_ID>
timeout = 100
worker_size = 10
```

It works with both ORG and GROUP level tokens. Be aware if you use group token, there's no confirmation messages on write operation. 


It is also possible to configure it using the cli
```bash
snykctl configure
token: 
group_id:
```

## Options
<table>
<tr><td>Flag</td><td>Description</td></tr>
<tr><td><kbd>-q</kbd></td><td>Quiet mode. Only shows ids</td></tr>
<tr><td><kbd>-n</kbd></td><td>Name only. Only shows names</td></tr>
<tr><td><kbd>-d</kbd></td><td>Debug mode. Prints HTTP requests</td></tr>
<tr><td><kbd>-t val</kbd></td><td>Timeout in seconds. </td></tr>
<tr><td><kbd>-p</kbd></td><td>Parallel flag. Use concurency if possible</td></tr>
<tr><td><kbd>-w val</kbd></td><td>Worker size (used with parallel)</td></tr>
<tr><td><kbd>-html</kbd></td><td>Html table</td></tr>
<tr><td><kbd>-lifecycle val</kbd></td><td>Lifecycle filter [ prod | dev | sandbox ]</td></tr>
<tr><td><kbd>-env</kbd></td><td>Environment filter [ front | back | onprem | mobile ]</td></tr>
</table>


