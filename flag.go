// import "github.com/u1f408/x"
package x

import (
    "os"
    "fmt"
    "flag"
)

type FlagSet struct {
    BaseName string
    Version string
    Usage func(*FlagSet)
    F *flag.FlagSet
}

func NewFlagSet(basename, version string) *FlagSet {
    n := flag.NewFlagSet(basename, flag.ExitOnError)
    f := &FlagSet {
        BaseName: basename,
        Version: version,
        F: n,
    }

    n.Usage = func() {
        f.Usage(f)
    }

    return f
}

func (f *FlagSet) Printf(format string, a ...any) (int, error) {
    return fmt.Fprintf(f.F.Output(), format + "\n", a...)
}

func (f *FlagSet) ParseFromOS() {
    f.Parse(os.Args[1:])
}

func (f *FlagSet) Parse(args []string) {
    err := f.F.Parse(args)
    if err == nil {
        return
    }

    if err == flag.ErrHelp {
        f.Usage(f)
        return
    }

    panic(err)
}

func (f *FlagSet) ExitUsage(code int) {
    f.Usage(f)
    os.Exit(code)
}
