package config

import (
	"github.com/imroc/req/v3"
)

type Config struct {
	Kubeconfig   string
	SecretPrefix string
	EnableDebug  bool
	Httpclient   *req.Client
}
