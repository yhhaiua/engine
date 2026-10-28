# engine

[![Go Reference](https://pkg.go.dev/badge/github.com/yhhaiua/engine.svg)](https://pkg.go.dev/github.com/yhhaiua/engine)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

基于 Go 的游戏网络框架，提供服务端/客户端长连接、粘包拆包、MySQL 定时落库、Redis/RabbitMQ、本地缓存、定时任务、日志和常用工具集。

## 特性

- **网络层**：TCP、WebSocket 服务端与客户端，统一 `net.Channel` 抽象，读写协程分离
- **协议编解码**：Netty 风格 `ByteBuf`，4 字节大端长度前缀自动拆包
- **数据持久化**：GORM + MySQL，自动建表，按 5 秒 / 30 秒 / 1 分钟 / 5 分钟分档批量落库，关服强制刷盘
- **本地缓存**：Guava 风格 LoadingCache，支持 LRU、SLRU、TinyLFU、过期策略和命中统计
- **中间件**：Redis 连接池（redigo）、RabbitMQ（amqp091-go）封装
- **任务调度**：cron 表达式定时任务 + ants 协程池
- **日志**：zap 多分类日志，按天轮转
- **其他**：YAML 配置、httprouter 风格 HTTP 路由、pprof、XML 配置表、JSON 工具

## 环境要求

- Go 1.24+
- MySQL 8.0+（推荐 8.4，默认 `caching_sha2_password`）
- Redis、RabbitMQ 按需启用

## 安装

```bash
go get github.com/yhhaiua/engine@v1.1.15
```

## 快速开始

### TCP 服务端

实现 `net.SocketListener`，在 `OnData` 中处理已拆好的包体：

```go
package main

import (
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/net"
)

type echoListener struct{}

func (echoListener) OnConnected(conn net.Channel)    {}
func (echoListener) OnDisconnected(conn net.Channel) {}

// OnData 收到的 msg 已经去掉 4 字节长度头，只包含包体
func (echoListener) OnData(conn net.Channel, msg *buffer.ByteBuf) {
	net.NewTcpSession(conn).Post(&echoPacket{data: msg.Bytes()})
}

type echoPacket struct {
	data []byte
}

func (p *echoPacket) Write(buf *buffer.ByteBuf) { buf.WriteNBytes(p.data) }
func (p *echoPacket) Read(buf *buffer.ByteBuf)  { p.data = buf.ReadNBytes() }
func (p *echoPacket) Copy() buffer.PacketCons   { return &echoPacket{data: p.data} }
func (p *echoPacket) CodeId() int               { return 1 }
func (p *echoPacket) Module() string            { return "echo" }
func (p *echoPacket) GetStage() int             { return 0 }

func main() {
	net.NewTCPServer(":9000", echoListener{}).Listen()
}
```

TCP 入站帧格式为 **4 字节大端长度 + 包体**，`net` 层会自动拆包并把包体交给 `OnData`。
出站时 `Channel.WriteAndFlush` 是裸写，不会加长度头；请使用 `net.NewTcpSession(conn).Post(packet)`，它会自动写入长度前缀。

### MySQL 定时落库

实体实现 `db.IEntity`；需要使用框架自增 id 时再实现 `db.ISince` 并返回 `true`：

```go
package main

import (
	"fmt"
	"time"

	"github.com/yhhaiua/engine/db"
)

type User struct {
	db.Model
	Name string
	Age  int32
}

func (u *User) TableName() string      { return "user" }
func (u *User) GetId() int64           { return u.ID }
func (u *User) SetId(id int64)         { u.ID = id }
func (u *User) GetCron() time.Duration { return db.PRE5SECOND }
func (u *User) IsMerger() bool         { return false }
func (u *User) GetSince() bool         { return true }

func main() {
	cfg := &db.MysqlConfig{
		Path:         "127.0.0.1:3306",
		Config:       "charset=utf8mb4&parseTime=True&loc=Local",
		Dbname:       "game",
		Username:     "root",
		Password:     "123456",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}

	// 连接、AutoMigrate、启动定时落库协程
	db.Init(cfg, &User{})
	defer db.ShutDown()

	u := &User{Name: "player", Age: 18}
	db.CreateAuto(u, 1, 1)      // 写入缓存队列，由后台按 GetCron() 批量落库
	time.Sleep(6 * time.Second) // 示例中等待一次落库完成

	var got User
	if db.First(u.ID, &got) {
		fmt.Println(got.Name)
	}
}
```

### YAML 配置

```yaml
mysql:
  path: 127.0.0.1:3306
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: game
  username: root
  password: 123456
  max-idle-conns: 10
  max-open-conns: 100

redis:
  path: 127.0.0.1:6379
  password: ""
  max-idle-conns: 10
  max-open-conns: 100
  db-num: 0

rabbit:
  path: 127.0.0.1:5672/
  userName: guest
  passWord: guest
  virtualHost: /
```

```go
var cfg struct {
	Mysql  db.MysqlConfig     `yaml:"mysql"`
	Redis  gredis.RedisConfig `yaml:"redis"`
	Rabbit mq.RabbitConfig    `yaml:"rabbit"`
}
if !conf.LoadYamlConfig("config.yaml", &cfg) {
	panic("load config failed")
}
```

### HTTP 路由

```go
router := grouter.New()
router.GET("/hello/:name", func(w http.ResponseWriter, r *http.Request, ps grouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
})
log.Fatal(http.ListenAndServe(":8080", router))
```

### 本地缓存

```go
c := cache.NewLoadingCache(
	func(k cache.Key) (cache.Value, error) {
		return loadUser(k) // 建议返回指针
	},
	cache.WithMaximumSize(1000),
	cache.WithExpireAfterAccess(30*time.Second),
	cache.WithRemovalListener(func(k cache.Key, v cache.Value) {
		// 淘汰回调
	}),
)
v, err := c.Get(int64(1001))
```

## 包一览

| 包 | 说明 |
| --- | --- |
| `net` | TCP / WebSocket 服务端与客户端，统一 `Channel` 抽象 |
| `buffer` | Netty 风格 `ByteBuf`，二进制协议读写 |
| `handler` | TCP / WebSocket 长度字段拆包解码 |
| `grouter` | httprouter 风格的高性能 HTTP 路由 |
| `db` | GORM/MySQL 访问、AutoMigrate、定时批量落库 |
| `cache` | 本地缓存（LRU / SLRU / TinyLFU） |
| `gredis` | Redis 连接池与常用命令封装 |
| `mq` | RabbitMQ 连接、交换机/队列声明、消费 |
| `cron` | cron 表达式定时任务 |
| `job` | 全局 ants 协程池 |
| `log` / `logzap` | 日志门面 / zap 多分类日志 |
| `conf` | YAML 配置加载 |
| `perf` | pprof 性能分析端口 |
| `jsonx` | json-iterator 封装 |
| `gxml` | XML 配置表加载（含 etree） |
| `util` | 原子数、集合、并发容器、表达式解析、HTTP 客户端等工具 |

`util` 下还包含：`cast`（类型转换）、`concurrent`（并发 Map/List）、`containers`（容器与排序）、`exp`（数学表达式）、`httputil`（HTTP 请求）、`md`（MD5/HMAC）、`params`（表单参数绑定）、`queue`（合并队列）、`treemap`（有序 Map）。

## 目录结构

```text
engine/
├── buffer/       ByteBuf 与二进制协议编解码
├── cache/        本地缓存
├── conf/         YAML 配置加载
├── cron/         定时任务
├── db/           MySQL/GORM 与定时落库
├── gredis/       Redis 连接池
├── grouter/      HTTP 路由
├── gxml/         XML 配置表加载
├── handler/      TCP/WS 拆包解码
├── job/          ants 协程池
├── jsonx/        JSON 封装
├── log/          日志门面
├── logzap/       zap 多分类日志
├── mq/           RabbitMQ 封装
├── net/          TCP/WebSocket 网络层
├── perf/         pprof
└── util/         通用工具
```

## 注意事项

- MySQL DSN 必须包含 `parseTime=True`，建议使用 `charset=utf8mb4&loc=Local`，否则 `time.Time` 字段无法正常扫描。
- 未标注 `size` 的 `string` 字段默认建成 `varchar(255)`；多个字符串列做联合索引时注意 InnoDB 3072 字节的索引键上限。
- `gredis` 的 `db-num` 当前实现固定使用 0，切换 DB 需要修改代码。
- `cache.Put` 的值建议传指针；`Count`、统计信息等能力需要类型断言到具体实现。
- `db.Init` 连接失败时只记录日志、不返回 error，调用方需要保证 MySQL 可用。

## License

本项目采用 [MIT License](LICENSE)。

`grouter`、`util/cast`、`gxml/etree`、`util/treemap`、`util/exp` 等目录包含来自第三方项目的代码，各自保留上游许可证，详见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
