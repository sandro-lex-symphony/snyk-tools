import requests
import argparse
import configparser
import os
import json

snyk_url = 'https://snyk.io/api/v1'


session = requests.Session()

def read_conf():
    home = os.path.expanduser('~')
    file = home + '/.snyk-users.conf'
    config = configparser.ConfigParser()
    config.read(file)
    return config

def configure():
    home = os.path.expanduser('~')
    file = home + '/.snyk-users.conf'
    config = configparser.ConfigParser()
    config.read(file)
        
    if 'group_token' in config['DEFAULT']:
        gt = config['DEFAULT']['group_token'][-4:]
        group_token = input(f"Group Token [...{gt}]: ")
    else: 
        group_token = input(f"Group Token []: ")
    if 'group_id' in config['DEFAULT']:
        gi = config['DEFAULT']['group_id'][-4:]
        group_id = input(f"Group id [...{gi}]: ")
    else:
        group_id = input("Group Id: ")
    if group_token != '':
        config['DEFAULT']['group_token'] = group_token
    if group_id != '':
        config['DEFAULT']['group_id'] =  group_id
    with open(file, 'w') as configfile:
        config.write(configfile)


def _get_orgs():
    conf = read_conf()
    url = snyk_url + '/orgs'
    headers = {'Authorization': 'token ' + conf['DEFAULT']['group_token']}
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    return res

def list_orgs():
    res = _get_orgs()
    if res.status_code != 200:
        return False

    access_rights = 0o755
    path = os.getcwd()
    cpath = path + '/symphony'
    os.mkdir(cpath, access_rights)
    for u in res.json()['orgs']:
        try:
            path = os.getcwd()
            opath = path + '/symphony/' + u['name']
            os.mkdir(opath, access_rights)
            with open(opath + '/meta.json', 'w') as outfile:
                json.dump(u, outfile)
        except OSError:
            print("error")


def _get_users(org):
    conf = read_conf()
    url = snyk_url + '/org/' + org + '/members'
    headers = {'Authorization': 'token ' + conf['DEFAULT']['group_token']}
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    return res

def get_users():
    context = get_context()
    res = _get_users(context['id'])
    path = os.getcwd()
    with open(path + '/users.json', 'w') as outfile:
        json.dump(res.json(), outfile)
    


def list_users(org, fmt='q'):
    res = _get_users(org)
    if fmt == 'json':
        print(res.json())
        return True
    elif fmt == 'all':
        for u in res.json():
            print(f"{u['id']} {u['role']} {u['username']} {u['name']} ")
    else:
        for u in res.json():
            print(f"{u['username']} {u['role']}")
        
    return True

def diff_users(src, dst):
    conf = read_conf()
    url = snyk_url + '/org/' + src + '/members'
    headers = {'Authorization': 'token ' + conf['DEFAULT']['group_token']}
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    users_src = res.json()

    url = snyk_url + '/org/' + dst + '/members'
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    users_dst = res.json()

    src_list = [ users_src[i]['name'] for i in range(len(users_src)) ]
    dst_list = [ users_dst[i]['name'] for i in range(len(users_dst)) ]
    
    column_2_size = 10
    column_1_size = len(src)

    # get the longest list
    if len(src) >= len(dst):
        l1 , l2 = sorted(src_list), sorted(dst_list)
    else:
        l1, l2 = sorted(dst_list), sorted(src_list)

    # create a set with all
    s1 = {l1[i] for i in range(len(l1))}
    s2 = {l2[i] for i in range(len(l2))}
    all = s1 | s2

    # print header
    out = src
    out += column_2_size * ' '
    out += dst + '\n'
    out += len(src) * '-'
    out += column_2_size * ' '
    out += len(dst) * '-'
    out += '\n'
    
    missing_str = '--> MISSING <--'
    for i in sorted(all):
        if i in l1:
            out += i
            filler_size = (column_1_size + column_2_size) - len(i)
        else:
            out += missing_str
            filler_size = (column_1_size + column_2_size) - len(missing_str)
        out += filler_size * ' '
        if i in l2:
            out += i
        else:
            out += missing_str
        out += '\n'

    print(out)

def copy_users(src, dst):
    # get list of users from org_src
    # for each user... 
    # POST add member to org_dst
    conf = read_conf()
    url = snyk_url + '/org/' + src + '/members'
    headers = {'Authorization': 'token ' + conf['DEFAULT']['group_token']}
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    users = res.json()
    headers['Content-Type'] = 'application/json'
    fail = False
    print(f"Going to copy users from {src} to {dst}")
    group_id = conf['DEFAULT']['group_id']
    for u in users:
        url = snyk_url + '/group/' + group_id + '/org/' + dst + '/members'
        payload = {'userId': u['id'], 'role': 'collaborator'}
        res = session.post(url, headers=headers, json=payload)
        if res.status_code != 200:
            print("Error: " + res.text)
            fail = True
        else:
            print("Added " + u['username'])
    
    return not fail

def get_orgs():
    conf = read_conf()
    url = snyk_url + '/orgs'
    headers = {'Authorization': 'token ' + conf['DEFAULT']['group_token']}
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    return res

def search_org(name):
    res = get_orgs()
    if res == False:
        return False
    out = ''
    for o in res.json()['orgs']:
        if name.lower() in o['name'].lower():
            out += o['id'] + '\t\t' + o['name'] + '\n'

    print(out)
    
def create_org(name):
    conf = read_conf()
    url = snyk_url + '/group/' + conf['DEFAULT']['group_id'] + '/org'
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token'], 
        'Content-Type': 'application/json'
    }

    payload = {'name': name}
    res = session.post(url, headers=headers, json=payload)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    print("ID: " + res.json()['id'])

def get_context():
    with open('meta.json') as json_file:
        return json.load(json_file)
         

def get_prjs():
    conf = read_conf()
    context = get_context()
    prjs = _search_projects(context['id'])
    path = os.getcwd()
    for p in prjs:
        print(p['name'])
        pname = p['name'].replace('/', '_')
        print(pname)
        os.mkdir(path + '/' + pname)
        with open(path + '/' + pname + '/meta.json', 'w') as outfile:
            json.dump(p, outfile)


def _search_projects(org):
    conf = read_conf()
    url = snyk_url + '/group/' + conf['DEFAULT']['group_id'] + '/org'
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token']
    }
    url = snyk_url + '/org/' + org + '/projects'
    res = session.get(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    return res.json()['projects']

def search_projects(org, fmt='simple', origin='', name='', delete=False):
    projects = _search_projects(org)
    if projects == False:
        return False

    # filters
    filtered = []
    for p in projects:
        if origin != '' and name != '':
            if p['origin'].lower() == origin.lower() and name.lower() in p['name'].lower():
                filtered.append(p)
        elif origin != '' and name == '':
            if p['origin'].lower() == origin.lower():
                filtered.append(p)
        elif origin == '' and name != '':
            if name.lower() in p['name'].lower():
                filtered.append(p)
        else:
            filtered.append(p)

    # no need to format    
    if delete:
        for p in filtered:
            delete_prj(org, p['id'])
        return True

    # format
    out = ''
    for p in filtered:
        if fmt == 'simple':
            out += p['id'] + '\t' + p['name'] + '\n'
        elif fmt == 'q':
            out += p['id'] + '\n'
        elif fmt == 'all':
            out += p['id'] + '\t' + p['created'] + '\t' + p['origin'] + '\t\t' + p['name'] + '\n'
    
    print(out)

def delete_prj(org, prj):
    conf = read_conf()
    url = snyk_url + '/group/' + conf['DEFAULT']['group_id'] + '/org'
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token']
    }
    url = snyk_url + '/org/' + org + '/project/' + prj
    res = session.delete(url, headers=headers)
    if res.status_code != 200:
        print(res.status_code)
        print("Error: " + res.text)
        return False
    print(prj + " DELETED")


def get_project_issues(org, prj):
    conf = read_conf()
    url = snyk_url + '/group/' + conf['DEFAULT']['group_id'] + '/org'
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token'],
        'Content-Type' : 'application/json'
    }
    url = snyk_url + '/org/' + org + '/project/' + prj + '/aggregated-issues'
    res = session.post(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    return res



def project_issues(org, prj, fmt='aggregated'):
    res = get_project_issues(org, prj)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    if fmt=='agregated':
        h = 0
        m = 0
        l = 0
        u = 0
        for i in res.json()['issues']:
            if i['issueData']['severity'] == 'high':
                h += 1
            elif i['issueData']['severity'] == 'medium':
                m += 1
            elif i['issueData']['severity'] == 'low':
                l += 1
            else:
                u += 1
        out = 'High: ' + str(h) + '\nMedium: ' + str(m) + '\nLow: ' + str(l)
        if u > 0:
            out += '\nUnknown: ' + str(u)
        print(out)
    else:
        out = ''
        for i in res.json()['issues']:
            ignored = ''
            if i['isIgnored'] == True:
                ignored = 'ignored'
            print(i['id'] + '\t' + i['pkgName'] + '\t' + i['issueData']['severity'] + '\t' + ignored)


def search_issue_org(org, issue):
    # get the list of prj_id 
    # for each prj ... 
    #    search for issue id
    
    projects = _search_projects(org)
    if projects == False:
        return False
    present = 0
    for p in projects:
        res = get_project_issues(org, p['id'])
        if res.status_code != 200:
            print("Error: " + res.text)
            return False
        
        for i in res.json()['issues']:
            if i['id'] == issue:
                print(p['id'] + '\t--' + p['name'] + ' PRESENT')
                present += 1
                continue
    if present == 0:
        print("Issue not found in the ORG")
        
    

def delete_ignore_issue(org, prj, issue):
    conf = read_conf()
    url = snyk_url + '/group/' + conf['DEFAULT']['group_id'] + '/org'
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token'],
        'Content-Type' : 'application/json'
    }
    url = snyk_url + '/org/' + org + '/project/' + prj + '/ignore/' + issue
    res = session.delete(url, headers=headers)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    print("OK")

def ignore_issue_group(issue):
    # get list of orgs
    # for each org
    #    get list of prjs
    #    for each prjs
    #        set ignore issue
    orgs = get_orgs()
    if orgs == False:
        return False
    for o in orgs.json()['orgs']:
        #projects = _search_projects(o['id'])
        #if projects == False:
        #    return False
        #qtde = len(projects)
        #print(o['name'] + ' ---> ' + str(qtde))
        search_issue_org(o['id'], issue)

    #print(orgs.json())


def get_depedencies(org, prj):
    url = 'org/orgId/dependencies'
    conf = read_conf()
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token'],
        'Content-Type' : 'application/json'
    }
    url = snyk_url + '/org/' + org + '/dependencies'
    payload = {
        'filters': {
            'projects': [ prj ]
            }
    }
    res = session.post(url, headers=headers, json=payload)
    if res.status_code != 200:
        print(res.status_code)
        print("Error: " + res.text)
        return False
    return res



def count_issues(org_id):
    total_projects, total_deps, total_issues, ratio = _count_issues(org_id, verbose=True)
    print(f"TOTAL: \t  {total_projects} \t {total_deps} \t {total_issues} \t  {ratio:0.6f}")    

def _count_issues(org_id, verbose=True):
    total_projects = 0
    total_issues = 0
    total_deps = 0
    th = 0
    tm = 0
    tl = 0
    tu = 0
    projects = _search_projects(org_id)
    for p in projects:
        res = get_project_issues(org_id, p['id'])
        if res == False:
            return False
        issues = len(res.json()['issues'])
        total_issues += issues
        h = 0
        m = 0
        l = 0
        u = 0
        for i in res.json()['issues']:
            if i['issueData']['severity'] == 'high':
                th += 1
            elif i['issueData']['severity'] == 'medium':
                tm += 1
            elif i['issueData']['severity'] == 'low':
                tl += 1
            else:
                tu += 1

        deps = get_depedencies(org_id, p['id'])
        deps = deps.json()['total']
        total_deps += deps
        if verbose:
            print(p['name'] + '\t' + str(deps)  + '\t' + str(issues))

    total_projects += len(projects)
    #ratio = 0
    #if total_deps != 0:
    #    ratio = total_issues / total_deps
    
    return total_projects, total_deps, total_issues, th, tm, tl, tu
    
def count_group_issues():
    orgs = _get_orgs()
    for o in orgs.json()['orgs']:
        tp, td, ti, th, tm, tl, tu = _count_issues(o['id'], verbose=False)
        if td != 0:
            r = ti / td
            rh = th / td
            rm = tm / td
            rl = tl / td
            tu = tu / td
        else:
            r = 0
            rh = 0
            rm = 0
            rl = 0 
            ru = 0
        print(f"{o['name']} \t\t {tp} \t {td} \t {ti} \t {th} \t {tm} \t {tl} \t {tu} \t {r:0.6f} \t {rh:0.6f} \t {rm:0.6f} \t {rl:0.6f} \t {rl:0.6f}")       

def _count_prjs():
    total = 0
    orgs = get_orgs()
    for o in orgs.json()['orgs']:
        projects = _search_projects(o['id'])
        print(o['name'] + '\t' + str(len(projects)))
        total += len(projects)

    print('TOTAL: \t' + str(total))    

def add_ignore_issue(org, prj, issue):
    conf = read_conf()
    url = snyk_url + '/group/' + conf['DEFAULT']['group_id'] + '/org'
    headers = {
        'Authorization': 'token ' + conf['DEFAULT']['group_token'],
        'Content-Type' : 'application/json'
    }
    url = snyk_url + '/org/' + org + '/project/' + prj + '/ignore/' + issue
    payload = {
        'ignorePath': '*',
        'reason': 'snyk-tool',
        'reasonType': 'not-vulnerable',
        'disregardIfFixable': False
    }
    res = session.post(url, headers=headers, json=payload)
    if res.status_code != 200:
        print("Error: " + res.text)
        return False
    print("OK")
    

