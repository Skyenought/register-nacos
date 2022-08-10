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

package resolver

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/client/discovery"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/registry-nacos/internal/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var _ discovery.Resolver = (*nacosResolver)(nil)

type (
	options struct {
		cluster string
		group   string
	}

	Option func(o *options)

	nacosResolver struct {
		client naming_client.INamingClient
		opts   options
	}
)

func (n *nacosResolver) Target(_ context.Context, target *discovery.TargetInfo) string {
	return target.Host
}

func (n *nacosResolver) Resolve(_ context.Context, desc string) (discovery.Result, error) {
	res, err := n.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: desc,
		HealthyOnly: true,
		GroupName:   n.opts.group,
		Clusters:    []string{n.opts.cluster},
	})
	if err != nil {
		return discovery.Result{}, err
	}
	if len(res) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %v", desc)
	}
	instances := make([]discovery.Instance, 0, len(res))
	for _, in := range res {
		if !in.Enable {
			continue
		}
		instances = append(instances,
			discovery.NewInstance("tcp", fmt.Sprintf("%s:%d", in.Ip, in.Port), int(in.Weight), in.Metadata))
	}
	if len(instances) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %v", desc)
	}
	return discovery.Result{
		CacheKey:  desc,
		Instances: instances,
	}, nil
}

func (n *nacosResolver) Name() string {
	return "nacos" + ":" + n.opts.cluster + ":" + n.opts.group
}

// WithCluster with cluster option.
func WithCluster(cluster string) Option {
	return func(o *options) { o.cluster = cluster }
}

// WithGroup with group option.
func WithGroup(group string) Option {
	return func(o *options) { o.group = group }
}

// NewDefaultNacosResolver create a default service resolver using nacos.
func NewDefaultNacosResolver(opts ...Option) discovery.Resolver {
	cli, err := nacos.NewDefaultNacosConfig()
	if err != nil {
		hlog.Errorf("Unexpected Error: %+v", err)
		return nil
	}
	return NewNacosResolver(cli, opts...)
}

// NewNacosResolver create a service resolver using nacos.
func NewNacosResolver(cli naming_client.INamingClient, opts ...Option) discovery.Resolver {
	op := options{
		cluster: "DEFAULT",
		group:   "DEFAULT_GROUP",
	}
	for _, option := range opts {
		option(&op)
	}
	return &nacosResolver{client: cli, opts: op}
}
