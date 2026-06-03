#!/bin/bash

ROUTER_CONTAINER_NAME='atg-nyc-bdr-01'
ROUTER_LABEL='ATG NYC BDR-01'

echo '🚀 Starting the default topology...'
./transitlab -start
echo '⏳ Topology started. BGP and routing state may take a minute to converge.'

echo "🔌 Connecting to ${ROUTER_LABEL}..."
docker exec -it "${ROUTER_CONTAINER_NAME}" bash -c "echo ''; echo '🔌 Connected to ${ROUTER_LABEL}.'; echo '🧭 Starting vtysh...'; echo '🧪 Try running `show ip route`!'; echo ''; vtysh"
