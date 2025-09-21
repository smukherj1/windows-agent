#!/usr/bin/bash

set -eu

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <IP_address>"
    echo "Error: An IP address must be provided as the first and only argument."
    exit 1
fi

ip="$1"
echo "Deploying to IP address: $ip"
scp out/agent.exe scripts/*.ps1 "$(whoami)@$ip:C:/Apps"