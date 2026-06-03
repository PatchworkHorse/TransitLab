## Default Topology

### 🤖 Current Autonomous Systems

#### StrataLink Telecom - AS65222

- IPv4: 172.40.64.0/20
- IPv6: 2001:db8:beef::/40

StrataLink Telecom is a regional ISP based in the Northeastern United States servicing residential, business, and wholesale customers. Currently their sole upstream provider is Axiom Global Transit.

#### Axiom Global Transit - AS65801

- IPv4: 100.100.100.0/18
- IPv6: 2001:db8:101::/36

Axiom Global Transit is a mid-sized Tier 2 transit provider with a global footprint. They peer with several Tier 1 networks and operate multiple data centers worldwide.

#### Core Nexus Exchange - AS65500

- IPv4: 198.51.100.0/23
- IPv6: 2001:db8:0:00::/64

Core Nexus Exchange is a neutral IX based in New York City. They provide a neutral switch fabric for participating networks to interconnect either bilaterally or multilateral via IX operated route servers.

### 🤝 Current Peering Arrangements & Upstreams

- **StrataLink Telecom (AS65222)** <--> **Axiom Global Transit (AS65801)**
  - Upstream transit relationship
  - Transit Subnets:
    - IPv4 100.100.101.0/30
