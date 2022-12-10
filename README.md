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
appmode = debug                  

# 日志文件路径 (默认为: "/var/log/pintd.log").
logfile = /var/log/pintd.log     

# 最大可打开文件数量
maxopenfiles = 8192
```

#### 重定向配置

```html
[redirect]

# 命名为test的重定向配置
[redirect.test]

# 协议, 可选TCP或UDP. (默认为 tcp)
proto = tcp                  

# pintd监听的IP地址. (默认 0.0.0.0，支持填写域名)
localaddr  = 0.0.0.0         

# pintd监听的端口. (默认 8888)
localport  = 8888     

# 重定向的IP地址. (默认 127.0.0.1，支持填写域名).
remoteaddr = 127.0.0.1    

# 重定向的端口. (默认 80).
remoteport = 80       

# 最大重定向的连接数 (默认 100. 此项仅对TCP有效).
maxredirects = 100

# 是否设置无延迟 NODELAY (仅对TCP有效)
nodelay = true

# 黑名单，列出的IP将被禁止访问pintd. (使用 ',' 分割IP地址).
denyaddrs = 123.34.77.0/24, 127.0.0.1

# 白名单, 将只允许列出的IP访问, 允许黑名单和白名单混合使用
admitaddrs = 221.0.0.0/8, 127.0.0.1
```

## 重载配置

重载配置不会影响已经建立的连接

```
pintd reload
```

## 注意

```
由于UDP没有连接的概念, 无法判断对端状态不能及时销毁协程、对应套接字等资源。
因此使用UDP协议时，请配合黑白名单限制IP地址，防止恶意的消耗资源。
```