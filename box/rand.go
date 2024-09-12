package main

import (
    "os"
    "fmt"
    "path"
    "strings"
    "math/rand"

    "github.com/u1f408/x"
)

var ToolRandAlphabets = map[string]string{
    "default": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
    "hex": "0123456789abcdef",
    "hexupper": "0123456789ABCDEF",
    "base32": "0123456789abcdefghjkmnpqrstvwxyz",
    "base32upper": "0123456789ABCDEFGHJKMNPQRSTVWXYZ",
    "upper": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
    "uppernum": "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
    "lower": "abcdefghijklmnopqrstuvwxyz",
    "lowernum": "abcdefghijklmnopqrstuvwxyz0123456789",
}

type ToolRand struct {
    Length *int
    Alphabet string
}

func (t *ToolRand) parseFlags(args []string) (*x.FlagSet, error) {
    fs := x.NewFlagSet(path.Base(os.Args[0]) + " rand", Version)
    fs.Usage = func(fs *x.FlagSet) {
        fs.Printf("usage: %s [options]", fs.BaseName)
        fs.Printf("alphabets:")
        for an, at := range ToolRandAlphabets {
            fs.Printf("  %s", an)
            fs.Printf("        %s", at)
        }

        fs.Printf("options:")
        fs.F.PrintDefaults()
    }

    t.Length = fs.F.Int("length", 24, "length of random text to generate")
    alphabetName := fs.F.String("alphabet", "default", "alphabet to use")
    fs.Parse(args)

    alphabet, ok := ToolRandAlphabets[*alphabetName]
    if !ok {
        return nil, fmt.Errorf("unknown alphabet")
    }

    t.Alphabet = alphabet
    return fs, nil
}

func (t *ToolRand) generate() string {
    sb := strings.Builder{}
    sb.Grow(*t.Length)
    for i := *t.Length - 1; i >= 0; i-- {
        sb.WriteByte(t.Alphabet[rand.Int63() % int64(len(t.Alphabet))])
    }

    return sb.String()
}

func (t *ToolRand) Description() string {
    return "generate random strings"
}

func (t *ToolRand) Run(args []string) error {
    _, err := t.parseFlags(args)
    if err != nil {
        return err
    }

    fmt.Println(t.generate())
    return nil
}
