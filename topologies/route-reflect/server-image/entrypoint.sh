#!/usr/bin/env bash
# Patchwork notes: 
# - This is becoming scoped to a route reflector, not a general server.
set -euxo pipefail

IFACE="${IFACE:-eth0}"
IPV4_ADDRESS="${IPV4_ADDRESS:-}"
IPV4_PREFIX="${IPV4_PREFIX:-24}"
IPV4_GATEWAY="${IPV4_GATEWAY:-}"

if [[ -z "${IPV4_ADDRESS}" ]]; then
  echo "IPV4_ADDRESS is required (example: 10.200.0.2 or 10.200.0.2/24)" >&2
  exit 1
fi

ip link set "${IFACE}" up
ip addr flush dev "${IFACE}"

if [[ "${IPV4_ADDRESS}" == */* ]]; then
  ip addr add "${IPV4_ADDRESS}" dev "${IFACE}"
else
  ip addr add "${IPV4_ADDRESS}/${IPV4_PREFIX}" dev "${IFACE}"
fi

if [[ -n "${IPV4_GATEWAY}" ]]; then
  ip route replace default via "${IPV4_GATEWAY}" dev "${IFACE}"
fi

# Start GoBGP daemon if config exists
GOBGP_CONFIG="${GOBGP_CONFIG:-/etc/gobgp/gobgpd.toml}"
GOBGP_API_HOSTS="${GOBGP_API_HOSTS:-:50051}"
if [[ -f "${GOBGP_CONFIG}" ]]; then
  echo "Starting GoBGP with config: ${GOBGP_CONFIG}"
  su -s /bin/bash -l debug -c "gobgpd -f ${GOBGP_CONFIG} --api-hosts ${GOBGP_API_HOSTS}" gobgp &
  sleep 1
else
  echo "Warning: GoBGP config not found at ${GOBGP_CONFIG}"
fi

exec "$@"
