package tekton

import (
	"strconv"
	"strings"

	"github.com/hongyuxuan/tekton-sdk-go/config"
	"github.com/hongyuxuan/tekton-sdk-go/core/errorx"
	"github.com/hongyuxuan/tekton-sdk-go/core/option"
	"github.com/hongyuxuan/tekton-sdk-go/service"
	v1 "github.com/hongyuxuan/tekton-sdk-go/service/v1"
	v1beta1 "github.com/hongyuxuan/tekton-sdk-go/service/v1beta1"
	"github.com/imroc/req/v3"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

type Client struct {
	Config *config.Config
	svcCtx *service.ServiceContext
}

func NewClient(opts ...option.ClientOptionFunc) *Client {
	config := &config.Config{}
	for _, opt := range opts {
		opt(config)
	}

	clientset, dynamicclient, token, baseUrl, err := createKubernetes(config)
	if err != nil {
		panic(err)
	}

	httpclient := req.C().
		OnBeforeRequest(func(client *req.Client, req *req.Request) error {
			if req.RetryAttempt > 0 {
				return nil
			}
			req.EnableDump()
			return nil
		}).
		OnAfterResponse(func(client *req.Client, res *req.Response) (err error) {
			responseCode := strconv.Itoa(res.StatusCode)
			if !strings.HasPrefix(responseCode, "2") && !strings.HasPrefix(responseCode, "3") {
				defer func() {
					if e := recover(); e != nil {
						err = errorx.NewError(int64(res.StatusCode), res.String(), nil)
					}
				}()
				resp := make(map[string]interface{})
				res.UnmarshalJson(&resp)
				err = errorx.NewError(int64(res.StatusCode), resp["message"].(string), resp)
			}
			if res.Err != nil {
				err = res.Err
			}
			return
		})
	httpclient.EnableInsecureSkipVerify().SetBaseURL(baseUrl)
	if config.EnableDebug {
		httpclient.EnableDebugLog()
		httpclient.EnableDumpAll()
	} else {
		httpclient.DisableDebugLog()
		httpclient.DisableDumpAll()
	}
	config.Httpclient = httpclient
	return &Client{
		Config: config,
		svcCtx: service.NewServiceContext(clientset, dynamicclient, config.SecretPrefix, token),
	}
}

func createKubernetes(c *config.Config) (clientset *kubernetes.Clientset, dynamicclient dynamic.Interface, token, baseUrl string, err error) {
	var conf *rest.Config
	if c.Kubeconfig != "" {
		conf, err = clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	} else {
		conf, err = rest.InClusterConfig()
		conf.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000) // setting a big ratelimiter for client-side throttling, default 5
	}
	if err != nil {
		return
	}
	clientset, err = kubernetes.NewForConfig(conf)
	if err != nil {
		return
	}
	dynamicclient, err = dynamic.NewForConfig(conf)
	if err != nil {
		return
	}
	token = conf.BearerToken
	baseUrl = conf.Host
	return
}

func (c *Client) Task(namespace string) *v1.Task {
	return v1.NewTask(c.Config, namespace, c.svcCtx)
}

func (c *Client) Pipeline(namespace string) *v1.Pipeline {
	return v1.NewPipeline(c.Config, namespace, c.svcCtx)
}

func (c *Client) PipelineRun(namespace string) *v1.PipelineRun {
	return v1.NewPipelineRun(c.Config, namespace, c.svcCtx)
}

func (c *Client) TriggerBinding(namespace string) *v1beta1.TriggerBinding {
	return v1beta1.NewTriggerBinding(c.Config, namespace, c.svcCtx)
}

func (c *Client) TriggerTemplate(namespace string) *v1beta1.TriggerTemplate {
	return v1beta1.NewTriggerTemplate(c.Config, namespace, c.svcCtx)
}

func (c *Client) EventListener(namespace string) *v1beta1.EventListener {
	return v1beta1.NewEventListener(c.Config, namespace, c.svcCtx)
}
