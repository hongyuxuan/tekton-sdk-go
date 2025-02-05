package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	tekton "github.com/hongyuxuan/tekton-sdk-go"
	"github.com/hongyuxuan/tekton-sdk-go/core/option"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SuiteTestClusterTask struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestClusterTask) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "hello"
	s.namespace = "default"
}

func (s *SuiteTestClusterTask) Test1CreateClusterTask() {
	yamlStr := `apiVersion: tekton.dev/v1beta1
kind: ClusterTask
metadata:
  name: hello
  labels:
    app: hello
spec:
  steps:
  - image: alpine:edge
    name: echo
    script: |
      #!/bin/sh
      echo "Hello World"
`
	err := s.client.ClusterTask("default").Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestClusterTask) Test2ListClusterTask() {
	res, err := s.client.ClusterTask(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=hello",
		Limit:         3,
	})
	s.Nil(err)
	if s.NotNil(res) {
		found := false
		for _, item := range res {
			if item.Name == s.name {
				b, _ := json.Marshal(item)
				fmt.Println(string(b))
				found = true
				break
			}
		}
		s.Equal(true, found)
	}
}

func (s *SuiteTestClusterTask) Test3GetClusterTask() {
	res, err := s.client.ClusterTask(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestClusterTask) Test4GetYamlClusterTask() {
	res, err := s.client.ClusterTask(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestClusterTask) Test5DeleteClusterTask() {
	err := s.client.ClusterTask(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestClusterTask(t *testing.T) {
	suite.Run(t, new(SuiteTestClusterTask))
}
