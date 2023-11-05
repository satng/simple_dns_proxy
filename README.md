# simple_dns_proxy 
Simple DNS forwarding,On the pc side we can modify the /etc/host file, but on the mobile side we need a simple DNS proxy

#### Custom domain name resolution
```json
{
  "domain_ip_map": {
    "www.123.com": "27.115.87.106",
    "example.net": "192.168.0.2"
  }
}

```

#### The  resolv.conf  file specifies the DNS servers that the system should use for name resolution
```text
nameserver 114.114.114.114
nameserver 223.5.5.5
nameserver 8.8.8.8
```
> 1.Just take the first value and adjust it according to your network
> 
> 2.If you do not want to set this parameter, you can use the default route address instead of eg.192.168.1.1