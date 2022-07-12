[中文README](README-cn)

---
# pintd

A tcp udp port redirector.

## Usage

```shell
pintd -h            Show Help Message.
pintd -c xxx        Specific Config File, Default pintd.ini.
```

## Configuation

Configuation file is ini format.

### Example

example config file : [example](pintd.ini)

#### global section. (Must Option)

```html
[pintd]        

# debug or release mode. (Must Option).                  
AppMode = debug         

# log file path. (default "/var/log/pintd.log").         
LogFile = /var/log/pintd.log     
```

#### redirect section

```html
[redirect]

# config section named 'test'.
[redirect.test]

# protocol, support tcp and udp. (Default tcp)
proto = tcp                  

# listening address of pintd. (Default 0.0.0.0)
localaddr  = 0.0.0.0         

# listening port of pintd. (Default 8888)
localport  = 8888     

# redirected address. (Default 127.0.0.1).
remoteaddr = 127.0.0.1    

# redirected port. (Default 80).
remoteport = 80       

# Maximum Redirect Connections, (Default 100. Only Valid for tcp).
maxredirects = 100    

# blacklist for deny some ip address to access pintd. (split ip address with ',').
denyaddrs = 123.34.77.0/24, 127.0.0.1
```
