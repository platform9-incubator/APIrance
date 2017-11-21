#!/usr/bin/env bash

function delete_all_functions()
{
    fission function list | grep -v NAME | awk '{print $1}'|xargs -I@ bash -c "fission function delete --name @"
}

function delete_all_routes()
{
    fission route list | grep -v NAME | awk '{print $1}'|xargs -I@ bash -c "fission route delete --name @"
}

function delete_all_envs()
{
    fission env list | grep -v NAME | awk '{print $1}'|xargs -I@ bash -c "fission env delete --name @"
}

delete_all_functions
delete_all_routes
delete_all_envs
