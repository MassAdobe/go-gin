/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: config
**/
package constants

import "time"

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:28 上午
 * @Description: 主进程常量
**/
const (
	SYSTEM_CONTROL_PWD = "pwd"
	CONFIG_NAME        = "/config.yml"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:02 下午
 * @Description: HTTP中的基本常量
**/
const (
	CONTENT_TYPE_KEY   = "Content-Type"                   // 请求协议种类键值
	CONTENT_TYPE_INNER = "application/json;charset=utf-8" // 请求协议种类内容
	REQUEST_USER_KEY   = "user"                           // 用户头信息键值
	REQUEST_TIMEOUT_TM = time.Duration(6 * time.Second)   // 五秒超时时间
	REQUEST_TRACE_ID   = "trace-id"                       // 跟踪ID
	REQUEST_STEP_ID    = "step-id"                        // 步骤ID
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:38 上午
 * @Description: 错误常量
**/
const (
	REQUEST_FINISH = "finish"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:39 上午
 * @Description: 过滤器常量
**/
const (
	HEADER_USER_JOSN          = `{"user_id": 1}` // 如果是开发环境，为了简化开发，默认给定user_id=1
	PROGRAM_ENV_DEBUG         = "debug"          // 开发环境
	RDS_IDEMPOTENT_HEADER_KEY = "idempotent"     // 头中幂等token参数
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:37 上午
 * @Description: 数据库常量
**/
const (
	DB_SQL      = "sql"
	DB_LOG      = "log"
	DB_TYPE     = "mysql"
	HAS_DELETED = "0" // 已删除
	NOT_DELETED = "1" // 未删除
	NOT_ENABLED = "0" // 无效
	IS_ENABLED  = "1" // 有效
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:40 上午
 * @Description: 上下文常量
**/
const (
	GO_CONTEXT_ENV_DEBUG      = "debug"
	CONTEXT_STATUS_FINISH     = "finish"
	CONTEXT_LOG_FUNCTION_MARK = "function"
	CONTEXT_LOG_PATH_NUM_MARK = "path_num"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:42 上午
 * @Description: 标点符号常量
**/
const (
	SPACE_MARK     = " "
	AND_MARK       = "&"
	COMMA_MARK     = ","
	QUESTION_MARK  = "?"
	FULL_STOP_MARK = "."
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:43 上午
 * @Description: token常量
**/
const (
	TOKEN_EXPIRE_TIME = 3600 // 幂等token有效时间(s)
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:44 上午
 * @Description: 日志常量
**/
const (
	TIME                 = "time"
	LOG_LEVEL            = "log_level"
	LOGGER               = "logger"
	DESC                 = "desc"
	MSG                  = "msg"
	TRACE                = "trace"
	ERROR                = "error"
	TIME_FORMAT          = "2006-01-02 15:04:05.000"
	ZAP_FIELD_TYPE       = "zapcore.Field"
	GIN_CONTEXT_TYPE     = "*gin.Context"
	GIN_CONTEXT_TYPES    = "[]*gin.Context"
	LOG_SERVER_NAME_MARK = "server_name"
	LOG_LEVEL_DEBUG      = "debug"
	LOG_LEVEL_INFO       = "info"
	LOG_LEVEL_ERROR      = "error"
	FUNCTION_MARK        = "function"
	PATH_NUM_MARK        = "path_num"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:46 上午
 * @Description: nacos常量
**/
const (
	BACKSLASH_MARK                = "/"
	NACOS_CONTEXT_PATH            = "/nacos"
	NACOS_LOG_DIR                 = "/tmp/nacos/log"
	NACOS_LOG_CACHE_DIR           = "/tmp/nacos/cache"
	NACOS_ROTATE_TIME             = "1h"
	NACOS_MAX_AGE                 = 3
	NACOS_LOG_LEVEL               = "debug"
	NACOS_SCHEMA                  = "http"
	NACOS_NOT_LOAD_CACHE_AT_START = true
	NACOS_SERVER_CONFIGS_MARK     = "serverConfigs"
	NACOS_CLIENT_CONFIG_MARK      = "clientConfig"
	LOG_LEVEL_MODIFIED_DEBUG      = "debug"
	LOG_LEVEL_MODIFIED_INFO       = "info"
	LOG_LEVEL_MODIFIED_WARN       = "warn"
	LOG_LEVEL_MODIFIED_ERROR      = "error"
	LOG_LEVEL_MODIFIED_DPANIC     = "dpanic"
	LOG_LEVEL_MODIFIED_PANIC      = "panic"
	LOG_LEVEL_MODIFIED_FATAL      = "fatal"
	NACOS_REGIST_IDC_MARK         = "idc"
	NACOS_REGIST_IDC_INNER        = "shanghai"
	NACOS_REGIST_TIMESTAMP_MARK   = "timestamp"
	NACOS_REGIST_VERSION_MARK     = "version"
	RATE_ALL                      = "all" // 限流漏斗全局
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 1:41 下午
 * @Description: redis的命令常量
**/
const (
	REDIS_DIAL_TCP = "tcp"
	REDIS_PONG     = "PONG"
	RDS_PING       = "PING"
	RDS_HEXISTS    = "hexists"
	RDS_HSET       = "hset"
	RDS_HGET       = "hget"
	RDS_HDEL       = "hdel"
	RDS_HGETALL    = "hgetall"
	RDS_HKEYS      = "hkeys"
	RDS_EXISTS     = "exists"
	RDS_GET        = "get"
	RDS_SET        = "set"
	RDS_ZADD       = "zadd"
	RDS_ZREVRANK   = "zrevrank"
	RDS_ZCOUNT     = "zcount"
	RDS_DEL        = "del"
	RDS_SETEX      = "setex"
	RDS_HLEN       = "hlen"
	RDS_INCR       = "incr"
	RDS_HINCRBY    = "hincrby"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:49 上午
 * @Description: 路由常量
**/
const (
	ROUTER_RELEASE = "release" // 生产
	ROUTER_TEST    = "test"    // 测试
	ROUTER_DEBUG   = "debug"   // 开发
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:49 上午
 * @Description: 协程常量
**/
const (
	GO_ROUTINE_MARK = "goroutine "
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/20 10:50 上午
 * @Description: 校验常量
**/
const (
	CHINESE_TYPE          = "zh"
	VALIDATE_COMMENT_MARK = "comment"
)
