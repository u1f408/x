package main

import (
    "fmt"
    "net"
    "context"

    "github.com/wzshiming/socks5"
)

func makeResolverSystem() *net.Resolver {
    return &net.Resolver{}
}

func makeResolverPlain(c *ProxyDnsConfig) *net.Resolver {
    if c.Enable == false || c.Host == "" {
        return makeResolverSystem()
    }

    return &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{}
            return d.DialContext(ctx, network, c.Host + ":" + c.Port)
        },
    }
}

func makeResolverSocks(socksHostPort string, c *ProxyDnsConfig) *net.Resolver {
    if c.Enable == false || c.Host == "" {
        return makeResolverSystem()
    }

    return &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d, err := socks5.NewDialer(socksHostPort)
            if err != nil {
                return nil, fmt.Errorf("failed connecting to SOCKS proxy: %s", err.Error())
            }

            return d.DialContext(ctx, network, c.Host + ":" + c.Port)
        },
    }
}

func resolveAddr(saddr string, resolver *net.Resolver) (string, error) {
    addrs, err := resolver.LookupIP(context.Background(), "ip4", saddr)
    if err != nil {
        return "", fmt.Errorf("failed DNS lookup for %s: %s", saddr, err.Error())
    }

    if len(addrs) == 0 {
        return "", fmt.Errorf("no DNS A records returned for %s", saddr)
    }

    return addrs[0].String(), nil
}
