# registry-nacos (*这是一个由社区驱动的项目*)

[English](README.md)

使用 **nacos** 作为 **Hertz** 的注册中心

##  这个项目应当如何使用?

### 服务端

**[registry-nacos/example/server/main.go](example/server/main.go)**

```go
import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	nacosRegistry "github.com/hertz-contrib/registry-nacos/registry"
)

func main() {
	addr := "127.0.0.1:8888"
	r := nacosRegistry.NewDefaultNacosRegistry()
	h := server.Default(
		server.WithHostPorts(addr),
		server.WithRegistry(r, &registry.Info{
			ServiceName: "hertz.test.demo",
			Addr:        utils.NewNetAddr("tcp", addr),
			Weight:      10,
			Tags:        nil,
		}),
	)
	// ...
	h.Spin()
}

```

### 客户端

**[registry-nacos/example/client/main.go](example/client/main.go)**

```go
import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/registry-nacos/resolver"
)

func main() {
	client, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	client.Use(sd.Discovery(resolver.NewDefaultNacosResolver()))
	// ...
}
```
## 如何运行示例 ?

### docker 运行 nacos-server 

- make prepare

```bash
make prepare
```

### run server

```go
go run server/main.go
```
### run client

```go
go run client/main.go
```

```go
2022/07/26 13:52:47.310617 main.go:46: [Info] code =200, body ={"ping":"pong2"}
2022/07/26 13:52:47.311019 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.311186 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.311318 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.311445 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.311585 main.go:46: [Info] code = 200, body ={"ping":"pong2"}
2022/07/26 13:52:47.311728 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.311858 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.311977 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
2022/07/26 13:52:47.312107 main.go:46: [Info] code = 200, body ={"ping":"pong1"}
```

## 自定义 Nacos Client 配置

### 服务端
```go
import (
	"context"
	"sync"
	
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	nacosRegistry "github.com/hertz-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
	}
	
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	
	addr := "127.0.0.1:8888"
	r := nacosRegistry.NewNacosRegistry(cli)
	h := server.Default(
		server.WithHostPorts(addr),
		server.WithRegistry(r, &registry.Info{
			ServiceName: "hertz.test.demo",
			Addr:        utils.NewNetAddr("tcp", addr),
			Weight:      10,
			Tags:        nil,
		}),
	)
	// ...
	h.Spin()
}

```

### 客户端
```go
import (
	"context"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	nacosRegistry "github.com/hertz-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
	}

	naocsCli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		}
  )
	if err != nil {
		panic(err)
	}
	r := nacos_demo.NewNacosResolver(naocsCli)
	cli.Use(sd.Discovery(r))
  // ...
} 

```

## 兼容性

Nacos 2.0 和 1.X 版本的 nacos-sdk-go 是完全兼容的，[详情](https://nacos.io/en-us/docs/2.0.0-compatibility.html)

