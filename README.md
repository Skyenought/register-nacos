# registry-nacos (*This is a community driven project*)

[中文](README_CN.md)

Nacos as service discovery for Hertz.

## How to use?

### Server

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

### Client

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
## How to run example?

### run docker

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
2022/07/26 13:52:47.310617 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311019 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311186 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311318 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311445 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311585 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311728 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311858 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.311977 main.go:46: [Info] code = 200, body ={"ping":"pong"}
2022/07/26 13:52:47.312107 main.go:46: [Info] code = 200, body ={"ping":"pong"}
```

## Custom Nacos Client Configuration

### Server
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

### Client
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
	//
}

```


## Compatibility
The server of Nacos2.0 is fully compatible with 1.X nacos-sdk-go. [see](https://nacos.io/en-us/docs/2.0.0-compatibility.html)
