# pintd

TCP/UDP 端口重定向工具。

## 使用方法

```shell
pintd -h            打印帮助信息。
pintd -c xxx        指定配置文件, 默认为: pintd.ini.
```

## 配置文件

配置文件使用 ini 格式.

### 示例

配置文件示例 : [example](pintd.ini)

#### 全局配置

```html
[pintd]             
# 选择debug或release模式           
AppMode = debug                  

# 日志文件路径 (默认为: "/var/log/pintd.log").
LogFile = /var/log/pintd.log     
```

#### 重定向配置

```html
[redirect]

# 命名为test的重定向配置
[redirect.test]

# 协议, 可选TCP或UDP. (默认为 tcp)
proto = tcp                  

# pintd监听的IP地址. (默认 0.0.0.0)
localaddr  = 0.0.0.0         

# pintd监听的端口. (默认 8888)
localport  = 8888     

# 重定向的IP地址. (默认 127.0.0.1).
remoteaddr = 127.0.0.1    

# 重定向的端口. (默认 80).
remoteport = 80       

# 最大重定向的连接数 (默认 100. 此项仅对TCP有效).
maxredirects = 100    

# 黑名单，列出的IP将被禁止访问pintd. (使用 ',' 分割IP地址).
denyaddrs = 123.34.77.0/24, 127.0.0.1
```
