# pintd

A tcp udp port redirector.

## Usage

```shell
pintd -h            Show Help Message.
pintd -c xxx        Specific Config File, Default pintd.ini.
```

## Config

Configuation file is ini format.

### Example

```html
[pintd]                          # global configuation. (Must Option).
AppMode = debug                  # debug or release.    (Must Option).
LogFile = /var/log/pintd.log     # log file path.       (Default /var/log/pintd.log).

# parent section
[redirect]                       # redirect configuation. (Must Option).

# child section named 'test'.
# For This Configuation, If you access port 8888 of the pintd server, pintd will redirect connection to 127.0.0.1:80.
[redirect.test]              # child section 'test' belong to parent section 'redirect'
proto = tcp                  # protocol, support tcp and udp. (Default tcp).
localaddr  = 0.0.0.0         # listening address of pintd. (your ip address, Default 0.0.0.0)
localport  = 8888            # listening port of pintd. (Default 8888).
remoteaddr = 127.0.0.1       # redirected address. (Default 127.0.0.1).
remoteport = 80              # redirected port. (Default 80).
maxredirects = 100           # Maximum Redirect Connections, (Default 100. Only Valid for tcp).
denyaddrs = 123.34.77.0/24, 127.0.0.1   # deny access to pintd's ip address. (split ip address with ',').

# child section named 'test2'.
[redirect.test2]
proto = udp                  # protocol udp.
localaddr  = 0.0.0.0
localport  = 9999
remoteaddr = 127.0.0.1
remoteport = 9090
```
