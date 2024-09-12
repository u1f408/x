package main

import (
    "os"
    "fmt"
    "path"

    "github.com/u1f408/x"
)

var (
    Version string

    flags = x.NewFlagSet(path.Base(os.Args[0]), Version)
    debugFlag = flags.F.Bool("log.debug", false, "enable debug logging")
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
    flags.Usage = func(f *x.FlagSet) {
        f.Printf("proxyssh, from github.com/u1f408/x")
        if Version != "" {
            f.Printf("version %s", Version)
        }

        f.Printf("\nusage: %s [options] <host> <port>", f.BaseName)
        f.Printf("invocation as SSH proxy: ssh -oProxyCommand='%s %%h %%p' user@host\n", "/path/to/" + f.BaseName)
        f.Printf("options:")
        f.F.PrintDefaults()
    }

    xdgConfig, ok := os.LookupEnv("XDG_CONFIG_HOME")
    if !ok {
        xdgConfig = path.Join(os.Getenv("HOME"), ".config")
    }

    dumpFlag := flags.F.Bool("dump", false, "dump proxy configuration")
    configPath := flags.F.String("config", path.Join(xdgConfig, "u1f408-x", "proxyssh.yml"), "path to configuration file")

    flags.ParseFromOS()

    config := new(Config)
    if err := ParseConfig(config, *configPath); err != nil {
        Logf("failed to load config: %s", err.Error())
        os.Exit(2)
    }

    if *dumpFlag {
        for i, pr := range config.Proxies {
            fmt.Printf("proxy %d:\n", i)
            fmt.Printf("  domain: %s\n", pr.Domain)

            if pr.ProxyUrl != "" {
                fmt.Printf("  proxy (in config): %s\n", pr.ProxyUrl)
            }

            if pr.Lookup.Consul.Enable {
                fmt.Printf("  proxy (Consul lookup): %s\n", pr.Lookup.Consul.ServiceAddr)
            }

            fmt.Println()
        }

        os.Exit(0)
    }

    if flags.F.NArg() < 2 {
        flags.ExitUsage(1)
    }

    host := flags.F.Arg(0)
    port := flags.F.Arg(1)

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
