# pintd

A tcp udp port redirector

## Config

Storage Configuations in pintd.ini.

### Example

```html
[pintd]
AppMode = debug                  # debug or release.
LogFile = /var/log/pintd.log     # log file path.

# parent section
[redirect]                       # must option, dont't delete.

# child section.
# If you access port 8888 of the pintd server, pintd will redirect connection to 127.0.0.1:80.
[redirect.test]              # child 'test' belong to parent 'redirect'
proto = tcp                  # protocol, support tcp and udp.
localaddr  = 0.0.0.0         # listen address (your ip address), like 127.0.0.1, 192.168.122.10 ....
localport  = 8888            # listen port
remoteaddr = 127.0.0.1       # redirect to address.
remoteport = 80              # redirect to port
maxredirects = 100           # Maximum Redirect Connections, default 100. (Only Valid for tcp).
denyaddrs = 123.34.77.0/24, 127.0.0.1   # deny access to pintd's ip address, split ip address with ','.

[redirect.test2]
proto = tcp
localaddr  = 0.0.0.0
localport  = 8080
remoteaddr = 127.0.0.1
remoteport = 22
denyaddrs = 127.0.0.1

[redirect.test3]
proto = udp
localaddr  = 0.0.0.0
localport  = 9999
remoteaddr = 127.0.0.1
remoteport = 9090

[redirect.test4]
proto = tcp
localaddr  = 0.0.0.0
localport  = 9999
remoteaddr = 127.0.0.1
remoteport = 9090
```
