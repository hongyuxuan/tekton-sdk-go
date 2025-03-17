package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	tekton "github.com/hongyuxuan/tekton-sdk-go"
	"github.com/hongyuxuan/tekton-sdk-go/core/option"
	"github.com/stretchr/testify/suite"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SuiteTestTask struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestTask) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "hello"
	s.namespace = "tektoncd-default"
}

func (s *SuiteTestTask) Test1CreateTask() {
	yamlStr := fmt.Sprintf(`apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
spec:
  steps:
  - image: alpine:edge
    name: echo
    script: |
      #!/bin/sh
      echo "Hello World"
`, s.name, s.namespace, s.name)
	err := s.client.Task(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestTask) Test2ListTask() {
	res, _, err := s.client.Task(s.namespace).List(context.TODO(), metav1.ListOptions{
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

func (s *SuiteTestTask) Test3GetTask() {
	res, err := s.client.Task(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestTask) Test4GetYamlTask() {
	res, err := s.client.Task(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestTask) Test5ListAllPipeline() {
	var res []tektonv1.Task
	conti := ""
	var err error
	for {
		res, conti, err = s.client.Task(s.namespace).List(context.TODO(), metav1.ListOptions{
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

func (s *SuiteTestTask) Test6DeleteTask() {
	err := s.client.Task(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestTask(t *testing.T) {
	suite.Run(t, new(SuiteTestTask))
}
