---
name: diagnose-network-latency
description: Use when users report slow response times, when SLO latency budgets are breached, or when investigating the network layer as a performance bottleneck
source: Brendan Gregg "Systems Performance" (2020, 2nd ed.); Google SRE Workbook latency chapter; RFC 6349 TCP throughput testing
tags: [networking, performance, debugging, observability]
verified: true
---

# Diagnose Network Latency

Systematically isolate the source of network latency across the application stack to identify whether the cause is DNS, TCP, routing, application, or infrastructure.

## Why This Is Best Practice

**Adopted by:** Google SRE methodology (USE/RED method); Netflix performance engineering; Cloudflare network diagnostics playbook
**Impact:** Network latency contributes to 47% of web performance issues (Cloudflare 2022); structured diagnosis reduces mean time to identify (MTTI) from hours to minutes
**Why best:** Latency has many sources; random investigation wastes time; systematic layer-by-layer elimination pinpoints root cause efficiently

Sources: Brendan Gregg "Systems Performance" 2nd ed. (2020); Google SRE Workbook Ch. 4; RFC 6349 (2011)

## Steps

1. **Quantify and characterize the latency** — Collect p50, p95, p99 latency for the affected request type over the past 24 hours. Determine: Is it consistent or intermittent? Affects all users or specific regions? Correlates with time of day, traffic volume, or deployments? This narrows the hypothesis space before any tooling.

2. **Isolate the network layer** — Use ping and traceroute to establish baseline RTT to the target: `ping -c 100 <host>` (watch for jitter and packet loss). `traceroute -n <host>` shows hop-by-hop latency. A high-latency hop that doesn't change in subsequent hops is a red herring — focus on the hop where latency first increases.

3. **Measure DNS resolution time** — DNS adds latency on every new connection: `dig +stats <hostname>` shows query time. `dig @8.8.8.8 +stats <hostname>` tests external resolver. TTL values below 60 seconds increase resolver round trips. Slow DNS (>50 ms) is common and often overlooked.

4. **Analyze TCP connection setup** — TCP handshake latency = 1 RTT. Capture with: `curl -w "%{time_namelookup} %{time_connect} %{time_starttransfer} %{time_total}\n" -o /dev/null -s <url>`. High `time_connect` vs `time_namelookup` indicates routing or firewall inspection latency, not application latency.

5. **Check for packet loss and retransmissions** — `ss -s` shows retransmit counts. `netstat -s | grep retransmit`. TCP retransmits cause 200-3000 ms latency spikes (RTO timer). Packet loss of 1% can cause 10% throughput loss on bulk transfers. Use `iperf3` to measure bandwidth and packet loss.

6. **Profile TLS handshake overhead** — TLS adds 1-2 RTT per new connection. `openssl s_client -connect <host>:443 -debug` shows handshake timing. TLS session resumption and HTTP/2 connection reuse eliminate per-request TLS overhead. Check if clients are reusing connections.

7. **Identify bandwidth saturation** — On the server: `sar -n DEV 1 10` shows interface utilization. `nload` or `iftop` shows real-time bandwidth. Interface saturation causes queuing latency that adds 10-100 ms. On cloud instances, check network credit exhaustion (T-series AWS, e2-micro GCP).

8. **Examine receive and transmit buffers** — Small TCP buffers limit throughput: `sysctl net.core.rmem_max net.core.wmem_max`. For high-bandwidth long-latency paths (BDP > 4 MB), default 256 KB buffers are the bottleneck. Apply RFC 6349 buffer sizing: BDP = bandwidth × RTT.

9. **Correlate with infrastructure metrics** — Check: CPU steal time (noisy neighbor on shared hosts), NIC driver errors (`ethtool -S <iface>`), VPC Flow Log drops (AWS: `REJECT` actions in flow logs), and load balancer error rates. Dropped packets in the hypervisor layer appear as jitter in the guest OS.

10. **Test end-to-end with synthetic monitoring** — Deploy probes from multiple geographic locations using tools like Blackbox Exporter, Synthetic Monitoring (Grafana Cloud), or Catchpoint. Reproduce from outside your network to distinguish client-side from server-side latency. Compare internal vs external measurements.

## Rules

- Always measure before tuning; kernel TCP parameter changes without measurement can worsen performance.
- Focus on p99 latency, not averages — averages hide the tail latency that degrades user experience.
- Confirm packet loss before investigating application code; 1% packet loss multiplies TCP retransmit latency.
- Retest after every change; latency has multiple contributing factors and fixing one may reveal another.

## Common Mistakes

- **Blaming the network before measuring** — application processing time (database queries, serialization) frequently masquerades as network latency; use distributed tracing to separate.
- **Ignoring the first-hop latency** — latency between client and first load balancer often exceeds datacenter-internal latency; measure from the client perspective, not server perspective.
- **Missing DNS caching failures** — containerized workloads (Kubernetes) resolve DNS on every connection unless configured otherwise; high DNS latency in k8s is a common root cause.
- **Measuring throughput instead of latency** — iperf3 bulk transfer throughput does not reveal request latency; they are different metrics measuring different things.

## When NOT to Use

- Latency caused by application-level issues (slow queries, inefficient algorithms) confirmed by profiling — escalate to application performance analysis
- WAN latency inherent to geographic distance — physics cannot be fixed by diagnosis; recommend CDN, edge caching, or user relocation
