package main

import (
    "os"
    "fmt"
    "path"
    "flag"
)

var (
    debugFlag = flag.Bool("log.debug", false, "enable debug logging")
)

func Logf(format string, a ...any) (int, error) {
    format = fmt.Sprintf("[%s] %s\n", path.Base(os.Args[0]), format)
    return fmt.Fprintf(os.Stderr, format, a...)
}

func Debugf(format string, a ...any) (int, error) {
    if *debugFlag {
        format = fmt.Sprintf("DEBUG: %s", format)
        return Logf(format, a...)
    }

    return 0, nil
}

func main() {
    flag.CommandLine.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), "proxyssh, from github.com/u1f408/x\n")
        fmt.Fprintf(flag.CommandLine.Output(), "\n")
        fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options] <host> <port>\n", path.Base(os.Args[0]))
        fmt.Fprintf(flag.CommandLine.Output(), "invocation as SSH proxy: ssh -oProxyCommand='%s %%h %%p' user@host\n", "/path/to/" + path.Base(os.Args[0]))
        fmt.Fprintf(flag.CommandLine.Output(), "\n")
        fmt.Fprintf(flag.CommandLine.Output(), "options:\n")
        flag.CommandLine.PrintDefaults()
    }

    xdgConfig, ok := os.LookupEnv("XDG_CONFIG_HOME")
    if !ok {
        xdgConfig = path.Join(os.Getenv("HOME"), ".config")
    }

    dumpFlag := flag.Bool("dump", false, "dump proxy configuration")
    configPath := flag.String("config", path.Join(xdgConfig, "u1f408-x", "proxyssh.yml"), "path to configuration file")

    flag.Parse()

    config := new(Config)
    if err := ParseConfig(config, *configPath); err != nil {
        Logf("failed to load config: %s", err.Error())
        os.Exit(2)
    }

    if *dumpFlag {
        for i, pr := range config.Proxies {
            fmt.Printf("proxy %d:\n", i)
            fmt.Printf("  domain: %s\n", pr.Domain)

            if pr.SocksUrl != "" {
                fmt.Printf("  proxy (in config): %s\n", pr.SocksUrl)
            }

            if pr.Lookup.Consul.Enable {
                fmt.Printf("  proxy (Consul lookup): %s\n", pr.Lookup.Consul.ServiceAddr)
            }

            if pr.Dns.Enable {
                fmt.Printf("  DNS lookups: proxied - %s:%s\n", pr.Dns.Host, pr.Dns.Port)
            } else {
                fmt.Printf("  DNS lookups: not proxied\n", pr.Dns.Host, pr.Dns.Port)
            }

            fmt.Println()
        }

        os.Exit(0)
    }

    if flag.NArg() == 0 {
        flag.CommandLine.Usage()
        os.Exit(1)
    }

    host := flag.Arg(0)
    port := flag.Arg(1)

    for idx, pr := range config.Proxies {
        t, err := TargetInfoFor(&pr, host, port)
        if err != nil {
            Debugf("proxy %d: err=%s", idx, err.Error())
        }

        if err == nil {
            Debugf("proxy %d: matched, via %s", idx, t.SocksUrl)
            err = t.ExecuteProxy()
            if err != nil {
                Logf("failed to connect to proxy %d: %s", idx, err.Error())
                os.Exit(3)
            }

            return
        }
    }

    Logf("couldn't find a suitable proxy for %s :(", host)
    os.Exit(1)
}
