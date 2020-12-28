# go-gin
### 基于gin框架的web脚手架

---

>#### web架构基于gin架构，增加了相关功能
>> + 优雅停服；
>> + 基于nacos的服务注册与发现（可选）；
>> + 基于nacos实现服务间调用的负载均衡（轮询）；
>> + 基于nacos的配置中心（可选）；
>> + 基于zap的日志管理（可选）；
>> + 动态压缩日志，打包，滚动；
>> + 修改了zap的日志输出为json格式，方便后期进行日志分析；
>> + 重写zap日志的打印方法，并重新定义日志输出引用所在文件、方法和代码行数；
>> + 增加了框架基于内网和nacos的服务注册发现的服务间调用，支持GET、POST、PUT、DELETED请求方法（可选）；
>> + 增加了内部调用重试机制，可在nacos中进行重试次数配置（可选）；
>> + 基于gorm框架，实现了读写分离双连接池，可在nacos中进行配置（可选）；
>> + 封装了GET、POST、PUT、DELETE的HTTP请求方法；
>> + 增加了面向前面的请求与响应日志输出；
>> + 增加了全局统一的错误处理（包括日志输出）;
>> + 增加了TRACE-ID和STEP-ID的日志打印（可选）；
>> + 增加了TRACE-ID和STEP-ID在服务间调用的应用和STEP-ID的自增；
>> + 增加了本地配置文件与代码解耦；
>> + 基于validated，完成RestFul接口参数校验；
>> + 动态调节日志输出级别（可选）；
>> + 基于nacos的自定义配置，动态更新（可选）；
>> + redis连接池（可选）；
>> + 增加了关闭系统时，释放Redis连接池与数据库连接池的句柄；

## 配置文件
```yaml
# 是否开启nacos配置中心
NacosConfiguration: true
# 是否开启nacos服务注册于发现
NacosDiscovery: true
# nacos地址
NacosServerIps: '127.0.0.1'
# nacos端口号
NacosServerPort: 8848
# nacos命名空间
NacosClientNamespaceId: 'f3e0c037-7fe1-452f-8f37-16b3810846b5'
# 请求Nacos服务端的超时时间（ms）
NacosClientTimeoutMs: 5000
# nacos配置文件名称
NacosDataId: 'go-framework.yml'
# nacos配置组名称
NacosGroup: 'go-framework'
# 日志输出路径(本地配置优先级最高)
LogPath: ''
# 日志级别(本地配置优先级最高)
LogLevel: ''
```

## nacos配置文件
```yaml
# 服务的配置
serve:
  # 端口号
  port: 28080
  # 服务名称
  server-name: 'go-framework'
  # 服务权重
  weight: 10
# 日志配置
log:
  # 日志地址
  path: '/Users/xxx/logs/golang'
  # 日志级别（可动态修改）
  level: 'info'
# 数据库配置（可选）
gorm:
  # 读库配置
  read: 
    # 数据库用户名
    username: 'root'
    # 数据库密码
    password: 'abcdefg'
    # 数据库IP
    ip: '127.0.0.1'
    # 数据库端口
    port: 3306
    # 数据库名称
    dbname: 'go_framework'
  # 写库配置
  write: 
    # 数据库用户名
    username: 'root'
    # 数据库密码
    password: 'abcdefg'
    # 数据库IP
    ip: '127.0.0.1'
    # 数据库端口
    port: 3306
    # 数据库名称
    dbname: 'go_framework'
# 内部调用配置
feign:
  # 内部调用重试次数（可选）
  retry-num: 3
# redis配置（可选）
redis:
  # redis的IP和端口
  ip-port: '127.0.0.1:6379'
  # 密码口令
  password: 'xxx'
  # 数据库号：0-15
  database: 0
  # 连接超时时间（秒）
  connect-timeout: 2
  # 读超时时间（秒）
  read-timeout: 2
  # 写超时时间（秒）
  write-timeout: 2
  # 最大挂起数
  max-idle: 6
  # 最大活跃数
  max-active: 4
  # 最大挂起时间（秒）
  idle-timeout: 60

```

>#### 未做功能
>> + 服务熔断；
>> + 服务降级；
>> + 业务报错适配器；
