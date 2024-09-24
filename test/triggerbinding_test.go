package main

import (
	"context"
	"fmt"
	"testing"

	tekton "github.com/hongyuxuan/tekton-sdk-go"
	"github.com/hongyuxuan/tekton-sdk-go/core/option"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SuiteTestTriggerBinding struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestTriggerBinding) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "testtriggerbinding"
	s.namespace = "default"
}

func (s *SuiteTestTriggerBinding) Test1CreateTriggerBinding() {
	yamlStr := `apiVersion: triggers.tekton.dev/v1beta1
kind: TriggerBinding
metadata:
  annotations:
    fiops/author: fanpengfei
  labels:
    app: testtriggerbinding
  creationTimestamp: "2024-03-05T09:15:14Z"
  name: testtriggerbinding
  namespace: default
spec:
  params:
  - name: gitrevision
    value: $(body.ref)
  - name: namespace
    value: default
  - name: gitrepositoryurl
    value: $(body.project.git_ssh_url)
  - name: application
    value: $(extensions.app_name)
`
	err := s.client.TriggerBinding(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestTriggerBinding) Test2ListTriggerBinding() {
	res, err := s.client.TriggerBinding(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=testtriggerbinding",
		Limit:         3,
	})
	s.Nil(err)
	if s.NotNil(res) {
		found := false
		for _, item := range res {
			if item.Name == s.name {
				found = true
				break
			}
		}
		s.Equal(true, found)
	}
}

func (s *SuiteTestTriggerBinding) Test3GetTriggerBinding() {
	res, err := s.client.TriggerBinding(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestTriggerBinding) Test4GetYamlTriggerBinding() {
	res, err := s.client.TriggerBinding(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestTriggerBinding) Test5DeleteTriggerBinding() {
	err := s.client.TriggerBinding(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestTriggerBinding(t *testing.T) {
	suite.Run(t, new(SuiteTestTriggerBinding))
}
