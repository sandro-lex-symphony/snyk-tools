#!/usr/bin/python3

import requests
import argparse
import configparser
import os

import snyklib

'''
snyk get orgs
snyk get users [--org]
snyk copy users [dst] [--src]   
snyk compare users  [dst] [--src] 
snyk get prjs [--org] 
snyk get issues [prj] [--org]
snyk get ignores [prj] [--org] | global 
snyk get issues [prj]
snyk get issues --org 

snyk create org [name]
snyl set ignore --org --prj [issue]
snyk set ignore --global [issue]

snyk set org [id]
snyk set prj [id]

snyk search org [name]
snyk search prj --org 

snyk delete prj [id] --org 
snyk delete ignore --org --prj [issue]

snyk get count issues --prj | --org | global 
snyk get count prjs --org | global 
'''


snyk_url = 'https://snyk.io/api/v1'

parser = argparse.ArgumentParser()
parser.add_argument("command", nargs='+', help="\nlist-users, copy-users, compare-users, configure, search-org, create-org, search-prj, delete-prj")
parser.add_argument("-f", "--format", help="format output")

session = requests.Session()

args = parser.parse_args()

def usage():
    print(f"USAGE: \n")

if len(args.command) != 2:
    usage()
    quit()


if args.command[1] == 'orgs':
    if args.command[0] == 'get':
        snyklib.list_orgs()
elif args.command[1] == 'org':
    if args.command[0] == 'create':
        print("create Orgs")
    elif args.command[0] == 'get':
        print("get Org")
    elif args.command[0] == 'search':
        print("search Org")
elif args.command[1] == 'prjs':
    if args.command[0] == 'get':
        snyklib.get_prjs()
elif args.command[1] == 'users':
    if args.command[0] == 'get':
        snyklib.get_users()    
elif args.command[1] == 'conf':
    if args.command[0] == 'set':
        snyklib.configure()

else:
    usage()
