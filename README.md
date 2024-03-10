# Dynamic DNS update server

This is a small http server which can issue DNS A record updates in GCP DNS. It can be used to configure PFSense as
custom DNS provider.

## Building

```shell
go build
```