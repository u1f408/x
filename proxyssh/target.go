package main

import (
    "io"
    "os"
    "fmt"
    "strings"

    "github.com/wzshiming/socks5"
)

type TargetInfo struct {
    SocksUrl string
    TargetResolvedHost string
    TargetPort string
}

func TargetInfoFor(proxy *ProxyConfig, host, port string) (*TargetInfo, error) {
    var err error

    if !strings.HasSuffix(host, proxy.Domain) {
        return nil, fmt.Errorf("non-matching domain")
    }

    resolvedHost := strings.Clone(host)
    if proxy.StripDomain {
        resolvedHost = strings.TrimSuffix(resolvedHost, proxy.Domain)
    }

    socksUrl, err := ProxyLookup(&proxy.Lookup, proxy.SocksUrl)
    if err != nil {
        return nil, fmt.Errorf("failed proxy lookup: %s", err.Error())
    }

    info := TargetInfo{
        SocksUrl: socksUrl,
        TargetResolvedHost: resolvedHost,
        TargetPort: port,
    }

    return &info, nil
}

func (t *TargetInfo) ExecuteProxy() error {
    dialer, err := socks5.NewDialer(t.SocksUrl)
    if err != nil {
        return fmt.Errorf("failed connecting to SOCKS proxy: %s", err.Error())
    }

    conn, err := dialer.Dial("tcp", t.TargetResolvedHost + ":" + t.TargetPort)
    if err != nil {
        return fmt.Errorf("failed connecting to target: %s", err.Error())
    }

    defer conn.Close()
    copyFn := func(src io.Reader, dst io.Writer) error {
        for {
            _, err := io.Copy(dst, src)
            if err != nil {
                return err
            }
        }
    }

    go copyFn(conn, os.Stdout)
    copyFn(os.Stdin, conn)

    return nil
}
