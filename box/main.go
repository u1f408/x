package main

import (
    "os"
    "fmt"
    "path"
    "time"
    "strings"
    "math/rand"

    "github.com/u1f408/x"
)

type Tool interface {
    Description() string
    Run(args []string) error
}

var (
    Version string

    flags = x.NewFlagSet(path.Base(os.Args[0]), Version)
    debugFlag = flags.F.Bool("log.debug", false, "enable debug logging")

    KnownTools = map[string]Tool{
        "rand": new(ToolRand),
    }
)

func main() {
    flags.Usage = func(f *x.FlagSet) {
        f.Printf("box, from github.com/u1f408/x")
        if Version != "" {
            f.Printf("version %s", Version)
        }

        f.Printf("\nusage: %s [global options] <subcommand> ...\n", f.BaseName)
        f.Printf("global options:")
        f.F.PrintDefaults()
        f.Printf("\nsubcommands:")
        for sn, si := range KnownTools {
            f.Printf("  %s", sn)
            desc := strings.Split(si.Description(), "\n")
            for _, dl := range desc {
                f.Printf("        " + dl)
            }
        }
    }

    flags.ParseFromOS()
    if flags.F.NArg() < 1 {
        flags.ExitUsage(1)
    }

    cmdName := flags.F.Arg(0)
    args := flags.F.Args()[1:]
    cmd, ok := KnownTools[cmdName]
    if !ok {
        fmt.Fprintf(os.Stderr, "error: unknown tool: %s\n", cmdName)
        fmt.Fprintf(os.Stderr, "see %s -help for available tools\n", path.Base(os.Args[0]))
        os.Exit(1)
    }

    rand.Seed(time.Now().UnixNano())
    err := cmd.Run(args)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error running tool %s: %s", cmdName, err.Error())
        os.Exit(2)
    }
}
