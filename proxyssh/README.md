# proxyssh

SSH `ProxyCommand` helper for SOCKS proxies.

## Configuration

Default configuration path: `$XDG_CONFIG_HOME/u1f408-x/proxyssh.yml`

Expects a top-level `proxies` key containing an array of proxy entries.
Each proxy entry has the following keys:

- `domain` (required, string) - domain suffix for this proxy, including leading `.`
- `strip_domain` (optional, string, default `false`) - whether to strip the defined domain suffix before connecting through the proxy
- `proxy_url` (required unless `lookup` is defined, string) - URL (including scheme) to the SOCKS proxy to use. use `socks5://` for local DNS resolution, `socks5h://` for proxied DNS resolution
- `lookup` (optional) - mechanisms for finding a SOCKS proxy to use, if not explicitly defined
    - `consul` - a Consul service DNS name to query

### Example

```yaml
proxies:
  - domain: ".fawn-vibe.ts.net"
    lookup:
      consul: "tailscale-proxy-soupnet.service.consul"
```
