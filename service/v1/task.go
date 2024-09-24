package v1

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/tekton-sdk-go/config"
	"github.com/hongyuxuan/tekton-sdk-go/service"
	"github.com/hongyuxuan/tekton-sdk-go/types"
	"github.com/imroc/req/v3"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Task struct {
	svcCtx     *service.ServiceContext
	httpclient *req.Client
	config     *config.Config
	namespace  string
	token      string
}

func NewTask(c *config.Config, namespace string, svcCtx *service.ServiceContext) *Task {
	token, err := svcCtx.GetBearerToken(namespace)
	if err != nil {
		panic(err)
	}
	return &Task{
		svcCtx:     svcCtx,
		httpclient: c.Httpclient,
		config:     c,
		namespace:  namespace,
		token:      token,
	}
}

type ListTaskResponse struct {
	ApiVersion string          `json:"apiVersion"`
	Items      []tektonv1.Task `json:"items"`
}

// https://apiserver.cluster.local:6443/apis/tekton.dev/v1/namespaces/default/tasks?labelSelector=app.kubernetes.io%2Fversion%3D0.3&limit=500
func (t *Task) List(ctx context.Context, opts metav1.ListOptions) (resp []tektonv1.Task, err error) {
	req := t.httpclient.Get(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/tasks", t.namespace)).SetBearerAuthToken(t.token)
	if opts.LabelSelector != "" {
		req.SetQueryParam("labelSelector", opts.LabelSelector)
	}
	if opts.FieldSelector != "" {
		req.SetQueryParam("fieldSelector", opts.FieldSelector)
	}
	if opts.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", opts.Limit))
	} else {
		req.SetQueryParam("limit", "500") // default 500
	}
	var res ListTaskResponse
	if err = req.SetSuccessResult(&res).Do(ctx).Err; err != nil {
		return
	}
	return res.Items, nil
}

// https://apiserver.cluster.local:6443/apis/tekton.dev/v1/namespaces/default/tasks/:name
func (t *Task) Get(ctx context.Context, name string) (resp tektonv1.Task, err error) {
	if err = t.httpclient.Get(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/tasks/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&resp).
		Do(ctx).Err; err != nil {
		return
	}
	return
}

func (t *Task) GetYaml(ctx context.Context, name string) (string, error) {
	var task types.TektonResource
	if err := t.httpclient.Get(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/tasks/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&task).
		Do(ctx).Err; err != nil {
		return "", err
	}
	delete(task.Metadata.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	task.Status = nil
	manifest, _ := yaml.Marshal(task)
	return string(manifest), nil
}

func (t *Task) Delete(ctx context.Context, name string) (err error) {
	return t.httpclient.Delete(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/tasks/%s", t.namespace, name)).SetBearerAuthToken(t.token).Do(ctx).Err
}

func (t *Task) Create(ctx context.Context, yamlStr string) (err error) {
	return t.svcCtx.ApplyYaml(ctx, t.namespace, yamlStr, "Task")
}
