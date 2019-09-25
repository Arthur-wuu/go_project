# Account

账户管理
---
> 提供登录，注册，找回密码，账户信息查询，权限验证(email, phone, captcha, ga)，JWT token颁发，token刷新, GA绑定，解绑，自选信息查询等相关功能

- 运行
```bash
IngramdeMacBook-Pro:api-account nprog$ # 获取依赖包
IngramdeMacBook-Pro:api-account nprog$ go get

IngramdeMacBook-Pro:api-account nprog$ # 编译
IngramdeMacBook-Pro:api-account nprog$ go build

IngramdeMacBook-Pro:api-account nprog$ #运行依赖
IngramdeMacBook-Pro:api-account nprog$ ls -lah
total 45192
drwxr-xr-x  8 nprog  staff   256B  5  2 14:48 .
drwxr-xr-x  5 nprog  staff   160B  5  2 14:49 ..
-rw-r--r--  1 nprog  staff   296B  5  2 14:48 Dockerfile
-rwxr-xr-x  1 nprog  staff    22M  5  2 14:48 api-account
-rw-r--r--  1 nprog  staff   976B  5  2 14:48 config.yml
drwxr-xr-x  4 nprog  staff   128B  5  2 14:48 locales
-rw-r--r--  1 nprog  staff   219B  5  2 14:48 supervisord.conf
drwxr-xr-x  6 nprog  staff   192B  5  2 14:48 templates

IngramdeMacBook-Pro:api-account nprog$ # 运行选项
IngramdeMacBook-Pro:api-account nprog$ ./api-account --help
Usage of ./api-account:
  -alsologtostderr
    	log to standard error as well as files
  -c string
    	conf file. (default "config.yml")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

- 配置
> 程序支持yaml和环境变量两种配置方式，环境变量配置优先级要高于yaml。
> 环境变量配置参考yaml配置文件，将父级key和子级key转换成大写并用下滑线(_)拼接起来作为环境变量名字即可。
```yaml
server:
  #ENV SERVER_PORT
  port: 8080
  #ENV SERVER_DEBUG
  debug: false
db: # 数据库配置
  host: 127.0.0.1
  port: 3306
  user: root
  password:
  db: account
  # 连接池最小连接数
  max_idle_conn: 10
  # 连接池最大连接数
  max_open_conn: 100
redis: # redis配置
  host: 127.0.0.1
  port: 6379
  password:
  db:
token: # 私钥，用于服务鉴权，请保证复杂度并定期更换
  secret: WHATISSCRET
  # 过期时间 单位为 minute
  expiration: 1440
ses: # aws发邮件服务
  region: us-east-1
  access_key_id: AKIAIGKZLCWQMJHK3L6A
  secret_key: 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
  sender: ingram.su@blockshine.com
sns: # aws 发短信服务
  region: us-east-1
  access_key_id: AKIAIGKZLCWQMJHK3L6A
  secret_key: 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
ip_find: # ip地址所在地查询apikey https://ipfind.co/
  auth: 26c27bb8-84fa-47f8-ac8d-1fccdc720669
path_limits: # 接口访问频率限制
#  - path: /login
#    method: post
#    limit: 20
#    # 单位为 second
#    time: 3600
#  - path: /login/ga
#    method: post
#    limit: 20
#    time: 600
```

```bash
IngramdeMacBook-Pro:api-account nprog$ # 应用端口
IngramdeMacBook-Pro:api-account nprog$ export SERVER_PORT=80
IngramdeMacBook-Pro:api-account nprog$ # 数据库主机
IngramdeMacBook-Pro:api-account nprog$ export DB_HOST=127.0.0.1
```