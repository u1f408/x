# box

as in "toolbox" - miscellaneous tools that don't belong anywhere else

```
usage: box <subcommand> ...
```

## rand

```
usage: box rand [options]

options:
  -alphabet string
        alphabet to use (default "default")
  -length int
        length of random text to generate (default 24)
```

alphabets:

* default - `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
* upper - `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
* uppernum - `ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
* lower - `abcdefghijklmnopqrstuvwxyz`
* lowernum - `abcdefghijklmnopqrstuvwxyz0123456789`
* hex - `0123456789abcdef`
* hexupper - `0123456789ABCDEF`
* base32 - `0123456789abcdefghjkmnpqrstvwxyz`
* base32upper - `0123456789ABCDEFGHJKMNPQRSTVWXYZ`
