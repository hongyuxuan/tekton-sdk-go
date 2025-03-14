package v1beta1

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/tekton-sdk-go/config"
	"github.com/hongyuxuan/tekton-sdk-go/service"
	"github.com/hongyuxuan/tekton-sdk-go/types"
	automatikotypes "github.com/hongyuxuan/tekton-sdk-go/types/automatiko"
	"github.com/imroc/req/v3"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApprovalTask struct {
	svcCtx     *service.ServiceContext
	httpclient *req.Client
	config     *config.Config
	namespace  string
	token      string
}

func NewApprovalTask(c *config.Config, namespace string, svcCtx *service.ServiceContext) *ApprovalTask {
	token, err := svcCtx.GetBearerToken(namespace)
	if err != nil {
		panic(err)
	}
	return &ApprovalTask{
		svcCtx:     svcCtx,
		httpclient: c.Httpclient,
		config:     c,
		namespace:  namespace,
		token:      token,
	}
}

type ListApprovalTaskResponse struct {
	ApiVersion string                         `json:"apiVersion"`
	Items      []automatikotypes.ApprovalTask `json:"items"`
}

// https://apiserver.cluster.local:6443/apis/tekton.automatiko.io/v1beta1/namespaces/default/approvaltasks?labelSelector=app.kubernetes.io%2Fversion%3D0.3&limit=500
func (t *ApprovalTask) List(ctx context.Context, opts metav1.ListOptions) (resp []automatikotypes.ApprovalTask, err error) {
	url := fmt.Sprintf("/apis/tekton.automatiko.io/v1beta1/namespaces/%s/approvaltasks", t.namespace)
	if t.namespace == "" {
		url = "/apis/tekton.automatiko.io/v1beta1/approvaltasks"
	}
	req := t.httpclient.Get(url).SetBearerAuthToken(t.token)
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
	var res ListApprovalTaskResponse
	if err = req.SetSuccessResult(&res).Do(ctx).Err; err != nil {
		return
	}
	return t.processItems(res.Items), nil
}

// https://apiserver.cluster.local:6443/apis/tekton.automatiko.io/v1beta1/namespaces/default/approvaltasks/:name
func (t *ApprovalTask) Get(ctx context.Context, name string) (resp automatikotypes.ApprovalTask, err error) {
	if err = t.httpclient.Get(fmt.Sprintf("/apis/tekton.automatiko.io/v1beta1/namespaces/%s/approvaltasks/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&resp).
		Do(ctx).Err; err != nil {
		return
	}
	return
}

func (t *ApprovalTask) GetYaml(ctx context.Context, name string) (string, error) {
	var approvaltask types.TektonResource
	if err := t.httpclient.Get(fmt.Sprintf("/apis/tekton.automatiko.io/v1beta1/namespaces/%s/approvaltasks/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&approvaltask).
		Do(ctx).Err; err != nil {
		return "", err
	}
	delete(approvaltask.Metadata.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	approvaltask.Status = nil
	manifest, _ := yaml.Marshal(approvaltask)
	return string(manifest), nil
}

func (t *ApprovalTask) Delete(ctx context.Context, name string) (err error) {
	return t.httpclient.Delete(fmt.Sprintf("/apis/tekton.automatiko.io/v1beta1/namespaces/%s/approvaltasks/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		Do(ctx).Err
}

func (t *ApprovalTask) Create(ctx context.Context, yamlStr string) (err error) {
	return t.svcCtx.ApplyYaml(ctx, "", yamlStr, "ApprovalTask")
}

func (t *ApprovalTask) processItems(items []automatikotypes.ApprovalTask) []automatikotypes.ApprovalTask {
	for i := range items {
		delete(items[i].ObjectMeta.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		items[i].ObjectMeta.ManagedFields = nil
	}
	return items
}
