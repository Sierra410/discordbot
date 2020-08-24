#!/bin/sh

servicename=discordbot
binf=/usr/local/bin/$servicename
server=root@172.22.10.100

remote_start() {
        ssh "$server" "systemctl start $servicename"
}

remote_stop() {
        ssh "$server" "systemctl stop $servicename"
}

deploy() {
        tmp=$(mktemp /tmp/httpserver_XXXXXXXX)

        set -e

        go build -o "$tmp"
        onexit() {
                rm "$tmp"
                remote_start
        }
        trap onexit EXIT

        scp "$tmp" "$server:${binf}_download"
        ssh "$server" "systemctl stop $servicename; mv '${binf}_download' '$binf'"
}

case $1 in
start)
        remote_start
        ;;
stop)
        remote_stop
        ;;
deploy)
        deploy
        ;;
*)
        echo "Usage: $0 start/stop/deploy"
        exit 1
        ;;
esac
