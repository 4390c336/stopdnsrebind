# coredns-rebind-protection

## Name

*stopdnsrebind* - Coredns plugin that implement `--stop-dns-rebind` from dnsmasq.

## Description

With `stopdnsrebind` enabled, users are able to block addresses from upstream nameservers which are in the private ranges.

The import order of this plugin matters, it is possible that it will not work depending on the import order

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
