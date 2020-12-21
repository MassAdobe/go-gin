# go-gin
### 基于gin框架的web脚手架

---

>#### web架构基于gin架构，增加了相关功能
>> + 优雅停服；
>> + 基于nacos的服务注册与发现；
>> + 基于nacos实现服务间调用的负载均衡（轮询）；
>> + 基于nacos的配置中心；
>> + 基于zap的日志管理；
>> + 动态压缩日志，打包，滚动；
>> + 修改了zap的日志输出为json格式，方便后期进行日志分析；
>> + 重新zap日志的打印方法，并重新定义日志输出引用所在文件、方法和代码行数；
>> + 增加了框架基于内网和nacos的服务注册发现的服务间调用，支持GET、POST、PUT、DELETED请求方法；
>> + 增加了内部调用重试机制，可在nacos中进行重试次数配置；
>> + 基于gorm框架，实现了读写分离双连接池，可在nacos中进行配置；
>> + 封装了GET、POST、PUT、DELETE的HTTP请求方法；
>> + 增加了面向前面的请求与响应日志输出；
>> + 增加了全局统一的错误处理（包括日志输出）;
>> + 增加了TRACE-ID和STEP-ID的日志打印；
>> + 增加了TRACE-ID和STEP-ID在服务间调用的应用和STEP-ID的自增；
>> + 增加了本地配置文件与代码解耦；

## 配置文件
```
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

>#### 未做功能
>> + 动态调节日志输出级别；
>> + 服务熔断；
>> + 服务降级；
>> + redis连接池；
>> + 业务报错适配器；
