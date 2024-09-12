package main

import (
    "fmt"
    "net"
    "strings"
    "context"
)

func resolveSRV(saddr string, resolver *net.Resolver) (string, error) {
    _, srvs, err := resolver.LookupSRV(context.Background(), "", "", saddr)
    if err != nil {
        return "", fmt.Errorf("failed DNS lookup for %s: %s", saddr, err.Error())
    }

    if len(srvs) == 0 {
        return "", fmt.Errorf("no DNS SRV records returned for %s", saddr)
    }

    return fmt.Sprintf("%s:%d", strings.TrimSuffix(srvs[0].Target, "."), srvs[0].Port), nil
}


func ProxyLookup(l *ProxyLookupConfig, fallback string) (string, error) {
    var potential string
    var err error

    if l.Consul.Enable {
        potential, err = resolveSRV(l.Consul.ServiceAddr, makeResolverPlain(&l.Consul.Dns))
        if err == nil && len(potential) > 0 {
            return "socks5h://" + potential, nil
        }
    }

    if len(fallback) > 0 {
        return fallback, nil
    }

    return "", fmt.Errorf("couldn't find a proxy service and no fallback was provided")
}
