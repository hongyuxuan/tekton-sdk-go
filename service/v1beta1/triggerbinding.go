package v1beta1

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/tekton-sdk-go/config"
	"github.com/hongyuxuan/tekton-sdk-go/service"
	"github.com/hongyuxuan/tekton-sdk-go/types"
	"github.com/imroc/req/v3"
	tektonv1beta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TriggerBinding struct {
	svcCtx     *service.ServiceContext
	httpclient *req.Client
	config     *config.Config
	namespace  string
	token      string
}

func NewTriggerBinding(c *config.Config, namespace string, svcCtx *service.ServiceContext) *TriggerBinding {
	token, err := svcCtx.GetBearerToken(namespace)
	if err != nil {
		panic(err)
	}
	return &TriggerBinding{
		svcCtx:     svcCtx,
		httpclient: c.Httpclient,
		config:     c,
		namespace:  namespace,
		token:      token,
	}
}

type ListTriggerBindingResponse struct {
	ApiVersion string                         `json:"apiVersion"`
	Items      []tektonv1beta1.TriggerBinding `json:"items"`
}

// https://apiserver.cluster.local:6443/apis/triggers.tekton.dev/v1beta1/namespaces/default/triggerbindings?labelSelector=app.kubernetes.io%2Fversion%3D0.3&limit=500
func (t *TriggerBinding) List(ctx context.Context, opts metav1.ListOptions) (resp []tektonv1beta1.TriggerBinding, err error) {
	req := t.httpclient.Get(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggerbindings", t.namespace)).SetBearerAuthToken(t.token)
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
	var res ListTriggerBindingResponse
	if err = req.SetSuccessResult(&res).Do(ctx).Err; err != nil {
		return
	}
	return t.processItems(res.Items), nil
}

// https://apiserver.cluster.local:6443/apis/triggers.tekton.dev/v1beta1/namespaces/default/triggerbindings/:name
func (t *TriggerBinding) Get(ctx context.Context, name string) (resp tektonv1beta1.TriggerBinding, err error) {
	if err = t.httpclient.Get(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggerbindings/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&resp).
		Do(ctx).Err; err != nil {
		return
	}
	return
}

func (t *TriggerBinding) GetYaml(ctx context.Context, name string) (string, error) {
	var task types.TektonResource
	if err := t.httpclient.Get(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggerbindings/%s", t.namespace, name)).
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

func (t *TriggerBinding) Delete(ctx context.Context, name string) (err error) {
	return t.httpclient.Delete(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggerbindings/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		Do(ctx).Err
}

func (t *TriggerBinding) Create(ctx context.Context, yamlStr string) (err error) {
	return t.svcCtx.ApplyYaml(ctx, t.namespace, yamlStr, "TriggerBinding")
}

func (t *TriggerBinding) processItems(items []tektonv1beta1.TriggerBinding) []tektonv1beta1.TriggerBinding {
	for i := range items {
		delete(items[i].ObjectMeta.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		items[i].ObjectMeta.ManagedFields = nil
	}
	return items
}
