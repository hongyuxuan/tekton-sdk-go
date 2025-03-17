package main

import (
	"context"
	"fmt"
	"testing"

	tekton "github.com/hongyuxuan/tekton-sdk-go"
	"github.com/hongyuxuan/tekton-sdk-go/core/option"
	"github.com/stretchr/testify/suite"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SuiteTestPipeline struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestPipeline) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "testpipeline"
	s.namespace = "default"
}

func (s *SuiteTestPipeline) Test1CreatePipeline() {
	yamlStr := fmt.Sprintf(`apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
spec:
  tasks:
    - name: hello-task
      taskSpec:
        steps:
          - name: echo
            image: alpine
            script: |
              #!/bin/sh
              echo "Hello, World!"`, s.name, s.namespace, s.name)
	err := s.client.Pipeline(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestPipeline) Test2ListPipeline() {
	res, _, err := s.client.Pipeline(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", s.name),
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

func (s *SuiteTestPipeline) Test3ListPipelineWithoutNamespace() {
	res, _, err := s.client.Pipeline("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", s.name),
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

func (s *SuiteTestPipeline) Test4GetPipeline() {
	res, err := s.client.Pipeline(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipeline) Test5GetYamlPipeline() {
	res, err := s.client.Pipeline(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipeline) Test6ListAllPipeline() {
	var res []tektonv1.Pipeline
	conti := ""
	var err error
	for {
		res, conti, err = s.client.Pipeline(s.namespace).List(context.TODO(), metav1.ListOptions{
			Limit:    1,
			Continue: conti,
		})
		s.Nil(err)
		if s.NotNil(res) {
			for _, item := range res {
				fmt.Println(item.GetName())
			}
		}
		if conti == "" {
			break
		}
	}
}

func (s *SuiteTestPipeline) Test7DeletePipeline() {
	err := s.client.Pipeline(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestPipeline(t *testing.T) {
	suite.Run(t, new(SuiteTestPipeline))
}
