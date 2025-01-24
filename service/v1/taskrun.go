package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hongyuxuan/tekton-sdk-go/config"
	"github.com/hongyuxuan/tekton-sdk-go/service"
	"github.com/imroc/req/v3"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskRun struct {
	svcCtx     *service.ServiceContext
	httpclient *req.Client
	config     *config.Config
	namespace  string
	token      string
}

func NewTaskRun(c *config.Config, namespace string, svcCtx *service.ServiceContext) *TaskRun {
	token, err := svcCtx.GetBearerToken(namespace)
	if err != nil {
		panic(err)
	}
	return &TaskRun{
		svcCtx:     svcCtx,
		httpclient: c.Httpclient,
		config:     c,
		namespace:  namespace,
		token:      token,
	}
}

type ListTaskRunResponse struct {
	ApiVersion string             `json:"apiVersion"`
	Items      []tektonv1.TaskRun `json:"items"`
}

type TaskRunList []tektonv1.TaskRun

func (r *TaskRunList) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *TaskRunList) ToJsonStringPretty() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

// https://apiserver.cluster.local:6443/apis/tekton.dev/v1/namespaces/default/tasks?labelSelector=app.kubernetes.io%2Fversion%3D0.3&limit=500
func (t *TaskRun) List(ctx context.Context, opts metav1.ListOptions) (resp TaskRunList, err error) {
	req := t.httpclient.Get(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/taskruns", t.namespace)).SetBearerAuthToken(t.token)
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
	var res ListTaskRunResponse
	if err = req.SetSuccessResult(&res).Do(ctx).Err; err != nil {
		return
	}

	return t.processItems(res.Items), nil
}

// https://apiserver.cluster.local:6443/apis/tekton.dev/v1/namespaces/default/tasks/:name
func (t *TaskRun) Get(ctx context.Context, name string) (resp tektonv1.TaskRun, err error) {
	if err = t.httpclient.Get(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/taskruns/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&resp).
		Do(ctx).Err; err != nil {
		return
	}
	resp.ObjectMeta.ManagedFields = nil
	return
}

func (t *TaskRun) GetStatus(ctx context.Context, name string) (string, error) {
	var taskrun types.TektonResource
	if err := t.httpclient.Get(fmt.Sprintf("/apis/tekton.dev/v1/namespaces/%s/taskruns/%s", t.namespace, name)).
					SetBearerAuthToken(t.token).
					SetSuccessResult(&taskrun).
					Do(ctx).Err; err != nil {
					return "", err
	}
	manifest, _ := yaml.Marshal(taskrun.Status)
	return string(manifest), nil
}

func (t *TaskRun) processItems(items []tektonv1.TaskRun) []tektonv1.TaskRun {
	for i := range items {
		delete(items[i].ObjectMeta.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		items[i].ObjectMeta.ManagedFields = nil
	}
	return items
}
