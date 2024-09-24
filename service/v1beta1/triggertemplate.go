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

type TriggerTemplate struct {
	svcCtx     *service.ServiceContext
	httpclient *req.Client
	config     *config.Config
	namespace  string
	token      string
}

func NewTriggerTemplate(c *config.Config, namespace string, svcCtx *service.ServiceContext) *TriggerTemplate {
	token, err := svcCtx.GetBearerToken(namespace)
	if err != nil {
		panic(err)
	}
	return &TriggerTemplate{
		svcCtx:     svcCtx,
		httpclient: c.Httpclient,
		config:     c,
		namespace:  namespace,
		token:      token,
	}
}

type ListTriggerTemplateResponse struct {
	ApiVersion string                          `json:"apiVersion"`
	Items      []tektonv1beta1.TriggerTemplate `json:"items"`
}

// https://apiserver.cluster.local:6443/apis/triggers.tekton.dev/v1beta1/namespaces/default/triggertemplates?labelSelector=app.kubernetes.io%2Fversion%3D0.3&limit=500
func (t *TriggerTemplate) List(ctx context.Context, opts metav1.ListOptions) (resp []tektonv1beta1.TriggerTemplate, err error) {
	req := t.httpclient.Get(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggertemplates", t.namespace)).SetBearerAuthToken(t.token)
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
	var res ListTriggerTemplateResponse
	if err = req.SetSuccessResult(&res).Do(ctx).Err; err != nil {
		return
	}
	return res.Items, nil
}

// https://apiserver.cluster.local:6443/apis/triggers.tekton.dev/v1beta1/namespaces/default/triggertemplates/:name
func (t *TriggerTemplate) Get(ctx context.Context, name string) (resp tektonv1beta1.TriggerTemplate, err error) {
	if err = t.httpclient.Get(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggertemplates/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		SetSuccessResult(&resp).
		Do(ctx).Err; err != nil {
		return
	}
	return
}

func (t *TriggerTemplate) GetYaml(ctx context.Context, name string) (string, error) {
	var task types.TektonResource
	if err := t.httpclient.Get(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggertemplates/%s", t.namespace, name)).
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

func (t *TriggerTemplate) Delete(ctx context.Context, name string) (err error) {
	return t.httpclient.Delete(fmt.Sprintf("/apis/triggers.tekton.dev/v1beta1/namespaces/%s/triggertemplates/%s", t.namespace, name)).
		SetBearerAuthToken(t.token).
		Do(ctx).Err
}

func (t *TriggerTemplate) Create(ctx context.Context, yamlStr string) (err error) {
	return t.svcCtx.ApplyYaml(ctx, t.namespace, yamlStr, "TriggerTemplate")
}
