package option

import (
	"github.com/hongyuxuan/tekton-sdk-go/config"
)

type ClientOptionFunc func(*config.Config)

//	func WithBearerToken(token string) ClientOptionFunc {
//		return func(c *config.Config) {
//			c.Token = token
//		}
//	}
func WithKubeconfig(kubeconfig string) ClientOptionFunc {
	return func(c *config.Config) {
		c.Kubeconfig = kubeconfig
	}
}

func WithSecretPrefix(secretPrefix string) ClientOptionFunc {
	return func(c *config.Config) {
		c.SecretPrefix = secretPrefix
	}
}

// func WithBaseUrl(baseUrl string) ClientOptionFunc {
// 	return func(c *config.Config) {
// 		c.BaseUrl = baseUrl
// 	}
// }

func WithDebug(enable bool) ClientOptionFunc {
	return func(c *config.Config) {
		c.EnableDebug = enable
	}
}
