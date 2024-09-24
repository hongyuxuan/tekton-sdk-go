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

type SuiteTestEventListener struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestEventListener) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "testeventlistener"
	s.namespace = "default"
}

func (s *SuiteTestEventListener) Test1CreateEventListener() {
	yamlStr := `apiVersion: triggers.tekton.dev/v1beta1
kind: EventListener
metadata:
  annotations:
    fiops/author: hongyuxuan
  labels:
    app: testeventlistener
  creationTimestamp: "2024-03-07T01:43:33Z"
  name: testeventlistener
  namespace: default
spec:
  namespaceSelector: {}
  resources: {}
  serviceAccountName: tekton-triggers-gitlab-sa
  triggers:
  - bindings:
    - kind: TriggerBinding
      ref: fiops-trigger-binding
    name: fiops-match
    template:
      ref: fiops-pipeline-template
`
	err := s.client.EventListener(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestEventListener) Test2ListEventListener() {
	res, err := s.client.EventListener(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=testeventlistener",
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

func (s *SuiteTestEventListener) Test3GetEventListener() {
	res, err := s.client.EventListener(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestEventListener) Test4GetYamlEventListener() {
	res, err := s.client.EventListener(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestEventListener) Test5DeleteEventListener() {
	err := s.client.EventListener(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestEventListener(t *testing.T) {
	suite.Run(t, new(SuiteTestEventListener))
}
