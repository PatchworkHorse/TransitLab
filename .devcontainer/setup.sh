#!/bin/bash

#       ,--,
#      ()   \ 
#       /    \
#     _/______\_
#    (__________)
#     /  /  \  \
#    /  /    \  \
#    `"`      `"`

echo '🚀 Starting the default topology...'
./transitlab -start
echo '⏳ The lab is starting now. BGP sessions and other routing state may take a minute to come up.'

echo '🔌 Connecting to the ATG BOS BDR-01 router...'
echo '🧭 Opening vtysh so you can inspect the running topology.'
docker exec -it atg-bos-bdr-01 bash -c "echo ''; echo '🔌 Connected to ATG BOS BDR-01.'; echo '🧭 Starting vtysh...'; echo ''; vtysh"
