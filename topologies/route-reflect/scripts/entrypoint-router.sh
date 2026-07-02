#!/bin/bash

case "$ROUTER_TYPE" in border|interior|route-server)
        ;;
    *)
        echo "ERROR: Invalid ROUTER_TYPE '$ROUTER_TYPE'. Exiting."
        exit 1
        ;;
esac

echo "Configuring virtual ${ROUTER_TYPE} router ${HOSTNAME} for ${ISP_NAME} (AS${ASN})"

# # Replace container daemons file with our template
# cp /scripts/daemons-template-${ROUTER_TYPE} /etc/frr/daemons
cp /scripts/daemons-template-default /etc/frr/daemons

# Give FRR config files to the FRR user and group
if id -u frr >/dev/null 2>&1; then
    chown -R frr:frr /etc/frr
fi
chmod 775 /etc/frr
find /etc/frr -maxdepth 1 -type f -exec chmod 664 {} \;

# # Copy in our custom FRR config
# cp /scripts/frr-${HOSTNAME}.conf /etc/frr/frr.conf

# Enable IP forwarding at the kernel level
sysctl -w net.ipv4.ip_forward=1

# Remove only Docker-assigned IP addresses while preserving interface state
for iface in $(ls /sys/class/net/ | grep -v lo); do
    # Remove IPv4 addresses but keep the interface up and preserve MAC
    ip -4 addr flush dev $iface
done

# Start FRR
/etc/init.d/frr start

# Keepalive
sleep infinity