# snykctl

**A cmd line tool for interacting with snyk API**

To get the list of Snyk Orgs:
```bash
snykctl list-orgs
6c5fdc1f-e675-4321-9bb0-bd7a22a34a51    Mobile-Sandbox
f6910fd7-43a3-4e20-8327-6b621b7746b3    JDC On Prem
ecd201fd-2bf1-4ef0-b4b6-2989010b5d47    CS
d56b0951-e4dd-4139-8556-1987f15757bd    FPT
f1fb9797-3132-43c7-861b-1925948c3128    Platform
dc04df24-17fd-42d4-8e80-5c2c5800ec53    Discovery
```

or search for a specific org
```bash
snykctl search-org rele
8b587d86-66b7-4947-b98c-0242de8b70cd    Release20.5.5
07dfd36f-9193-40c2-8e44-2d5970231baa    release20.7
de6ec69b-5528-4e2b-bf20-14f6293cb273    release20.9
a583250c-36fa-40fe-b71a-abb8e3180993    Release20.9.1
ed4a4d1b-cc26-46ba-9142-07f88b6a50f4    Release20.12
c5807011-222a-42c5-bf1f-8660ebfb85a6    ReleaseTest
6e580461-1720-47da-90e9-7c94b6c0b96f    Release20.10-MT
4c961058-b36e-4510-894f-99d3f39d3498    Release20.13
4c986767-9543-43b2-85b9-2747a261fe91    Release 20.10
711c53b6-a85d-4a51-a34f-42552cc8572e    Release - Current
```

list projects for a project
```bash
snykctl list-projects 8b587d86-66b7-4947-b98c-0242de8b70cd
3802e440-e69c-4387-a3a3-9c0a4f2f69fa    com.symphony.sbe.core:malware-scan-client
8c995d18-5de4-4cb7-90fb-4de229f73be6    com.symphony.sbe.core:cache
3459efe9-7c3a-45fd-9138-d62661e572bf    com.symphony.sbe.core:retentioncommon
eac77497-4c16-4db0-94f6-0fa8c8e5ab91    com.symphony.sbe.core:data-common
404f5b82-2289-47f9-aecb-8b180371cce0    symphony-login
eb730354-5951-45d1-badb-e5d9b4fdd122    com.symphony.sbe.core:jwt-commons
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
* golang > 1.13
* Makefile

```bash
make
cp bin/snykctl /usr/local/bin/
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


