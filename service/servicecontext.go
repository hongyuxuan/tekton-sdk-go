package service

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/hongyuxuan/tekton-sdk-go/core/errorx"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	uyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
)

type ServiceContext struct {
	Clientset     *kubernetes.Clientset
	Dynamicclient dynamic.Interface
	SecretPrefix  string
	BearerToken   string
}

func NewServiceContext(clientset *kubernetes.Clientset, dynamicclient dynamic.Interface, secretPrefix, token string) *ServiceContext {
	return &ServiceContext{
		Clientset:     clientset,
		Dynamicclient: dynamicclient,
		SecretPrefix:  secretPrefix,
		BearerToken:   token,
	}
}

func (s *ServiceContext) GetBearerToken(namespace string) (token string, err error) {
	if s.BearerToken != "" {
		return s.BearerToken, nil
	}
	var res *corev1.SecretList
	res, err = s.Clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}
	secret, ok := lo.Find(res.Items, func(item corev1.Secret) bool {
		return strings.HasPrefix(item.Name, s.SecretPrefix)
	})
	if !ok {
		return "", errorx.NewDefaultError("cannot find secret with prefix=%s", s.SecretPrefix)
	}
	token = string(secret.Data["token"])
	return
}

func (s *ServiceContext) ApplyYaml(ctx context.Context, namespace, yamlStr, kind string) (err error) {
	d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(yamlStr), 4096)
	var unstructureObj *unstructured.Unstructured
	for {
		unstructureObj, err = s.getUnstructured(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if unstructureObj.GetKind() != kind {
			return errorx.NewDefaultError("Kind %s mismatch with %s", unstructureObj.GetKind(), kind)
		}
		var gvr schema.GroupVersionResource
		gvr, err = s.gtGVR(unstructureObj.GroupVersionKind())
		if err != nil {
			return
		}
		_, getErr := s.Dynamicclient.Resource(gvr).Namespace(namespace).Get(ctx, unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			_, createErr := s.Dynamicclient.Resource(gvr).Namespace(namespace).Create(ctx, unstructureObj, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
			return
		}

		if namespace == unstructureObj.GetNamespace() {
			_, err = s.Dynamicclient.Resource(gvr).Namespace(namespace).Update(ctx, unstructureObj, metav1.UpdateOptions{})
			if err != nil {
				return errorx.NewDefaultError("unable to apply yaml of resource[%s]: %s", unstructureObj.GetName(), err.Error())
			}
		} else {
			_, err = s.Dynamicclient.Resource(gvr).Update(ctx, unstructureObj, metav1.UpdateOptions{})
			if err != nil {
				return errorx.NewDefaultError("ns is nil unable to update resource: %s", err.Error())
			}
		}
	}
	return
}

func (s *ServiceContext) getUnstructured(d *uyaml.YAMLOrJSONDecoder) (unstructureObj *unstructured.Unstructured, err error) {
	var rawObj runtime.RawExtension
	err = d.Decode(&rawObj)
	if err == io.EOF {
		return
	}
	if err != nil {
		err = errorx.NewDefaultError("decode is err: %v", err.Error())
		return
	}
	obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
	if err != nil {
		err = errorx.NewDefaultError("rawobj is err: %v", err.Error())
		return
	}
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		err = errorx.NewDefaultError("tounstructured is err %v", err.Error())
		return
	}
	unstructureObj = &unstructured.Unstructured{Object: unstructuredMap}
	return
}

func (s *ServiceContext) gtGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	gr, err := restmapper.GetAPIGroupResources(s.Clientset.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}
