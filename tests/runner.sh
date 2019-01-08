#!/bin/bash -e

if [ -z "$1" ]
  then
    echo "No argument supplied"
    echo "Possible args:  up, down, kill, wait"
    exit 1
fi

type docker
type docker-compose
currdir="$(dirname $(realpath "$0"))"

function docker_compose_action()
{
    pushd "$currdir"/data/zabbix3
    docker-compose $1 $2
    popd
    pushd "$currdir"/data/chef
    docker-compose $1 $2
    popd
}

case $1 in
  up)
    docker_compose_action up -d
    ;;
  wait)
    while ! curl -k https://127.0.0.1:3443 --fail; do sleep 1; done
    docker exec -it chef-server chef-server-wait-lock
    ;;
  *)
    docker_compose_action $1
    ;;
esac
