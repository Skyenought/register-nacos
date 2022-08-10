// Copyright 2021 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const (
	nacosNamespaceID       = "public"
	nacosDefaultServerAddr = "127.0.0.1"
	nacosDefaultPort       = 8848
	nacosDefaultRegionID   = "cn-hangzhou"
)

// NewDefaultNacosConfig create a default Nacos client
func NewDefaultNacosConfig() (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacosDefaultServerAddr, nacosDefaultPort),
	}
	cc := constant.ClientConfig{
		NamespaceId:         nacosNamespaceID,
		RegionId:            nacosDefaultRegionID,
		NotLoadCacheAtStart: true,
		CustomLogger:        NewCustomNacosLogger(),
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
