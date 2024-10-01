# coredns-rebind-protection

## Name

*stopdnsrebind* - Coredns plugin that implement `--stop-dns-rebind` from dnsmasq.

## Description

With `stopdnsrebind` enabled, users are able to block addresses from upstream nameservers which are in the private ranges.

The import order of this plugin matters, it is possible that it will not work depending on the import order

# 🌐 Default Blocking Rules

### 🔒 Loopback Addresses
- **`127.0.0.1/8`**

### 🔒 Private Addresses
- **`10.0.0.0/8`**
- **`172.16.0.0/12`**
- **`192.168.0.0/16`**

### 🔒 Link Local Addresses
- **`169.254.0.0/16`**

### 🔒 Unspecified
- **`0.0.0.0`**

### 🔒 Interface Local Multicast
- **`224.0.0.0/24`**

### 🔒 DenyList
- Add your entries in the plugin configuration.

---

Keeping the network secure! 🔐

## Syntax

```
stopdnsrebind [ZONES...] {
    allow [ZONES...]
    deny [IPNet]
}
```

- **ZONES** zones that are allowed o resolve to private addresses

## Examples

To demonstrate the usage of plugin stopdnsrebind, here we provide some typical examples.

~~~ corefile
. {
    stopdnsrebind {
        allow internal.example.org
        deny 192.0.2.1/24
    }
}
~~~
