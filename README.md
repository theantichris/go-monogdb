# podcasts

Application for tracking podcast information using Go and MongoDB Atlas.

## Connection issues

If the application is unable to connect to Atlas add the following line to the top of `/etc/resolv.conf`:

```shell
nameserver 8.8.838
```
