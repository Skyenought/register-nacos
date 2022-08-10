package registry_nacos

import "github.com/nacos-group/nacos-sdk-go/clients/naming_client"

type (
	options struct {
		cluster string
		group   string
	}

	// Option is nacos option.
	Option func(o *options)

	nacosRegistry struct {
		client naming_client.INamingClient
		opts   options
	}

	nacosResolver struct {
		client naming_client.INamingClient
		opts   options
	}
)

// WithCluster with cluster option.
func WithCluster(cluster string) Option {
	return func(o *options) {
		o.cluster = cluster
	}
}

// WithGroup with group option.
func WithGroup(group string) Option {
	return func(o *options) {
		o.group = group
	}
}
