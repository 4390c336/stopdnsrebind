# stopdnsrebind

## Name

*stopdnsrebind* - Coredns plugin that implement `--stop-dns-rebind` from dnsmasq.

## Description

With `stopdnsrebind` enabled, users are able to block addresses from upstream nameservers which are in the private ranges.

## Syntax

```
stopdnsrebind [ZONES...] {
    allow [ZONES...]
}
```

- **ZONES** zones that are allowed o resolve to private addresses

## Examples

To demonstrate the usage of plugin stopdnsrebind, here we provide some typical examples.

~~~ corefile
. {
    stopdnsrebind {
        allow internal.example.org
    }
}
~~~
