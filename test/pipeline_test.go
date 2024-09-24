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
	yamlStr := `apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  annotations:
    fiops/author: hongyuxuan
  creationTimestamp: "2024-04-08T03:35:44Z"
  name: testpipeline
  namespace: default
  labels:
    app: testpipeline
spec:
  params:
  - name: revision
    type: string
  - default: git@gitlab.cicconline.com:xficc/devops/fiops.git
    name: repo-url
    type: string
  tasks:
  - name: git-clone
    params:
    - name: url
      value: $(params.repo-url)
    - name: subdirectory
      value: ""
    - name: deleteExisting
      value: "true"
    - name: revision
      value: $(params.revision)
    taskRef:
      kind: Task
      name: git-clone
    workspaces:
    - name: output
      workspace: shared-workspace
    - name: ssh-directory
      workspace: git-credentials
  - name: build
    params:
    - name: workbase
      value: feishu/api
    - name: packages
      value: feishuService
    - name: flags
      value: -o
    runAfter:
    - git-clone
    taskRef:
      kind: Task
      name: go1.19-build
    workspaces:
    - name: source
      workspace: shared-workspace
  - name: docker-build
    params:
    - name: IMAGE
      value: repo.cicc.com.cn/fi-fiops-docker-local/devops/feishu:$(tasks.build.results.version)-$(tasks.build.results.timestamp)-$(tasks.git-clone.results.commit)
    - name: CONTEXT
      value: feishu/api
    runAfter:
    - build
    taskRef:
      kind: Task
      name: kaniko
    workspaces:
    - name: source
      workspace: shared-workspace
    - name: dockerconfig
      workspace: dockerhub-auth
  - name: deploy
    params:
    - name: image-url
      value: $(tasks.docker-build.results.IMAGE_URL)
    - name: app-name
      value: feishu
    runAfter:
    - docker-build
    taskRef:
      kind: Task
      name: k8s-deploy
  workspaces:
  - name: shared-workspace
  - name: dockerhub-auth
  - name: git-credentials
`
	err := s.client.Pipeline(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestPipeline) Test2ListPipeline() {
	res, err := s.client.Pipeline(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=testpipeline",
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

func (s *SuiteTestPipeline) Test3GetPipeline() {
	res, err := s.client.Pipeline(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipeline) Test4GetYamlPipeline() {
	res, err := s.client.Pipeline(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipeline) Test5DeletePipeline() {
	err := s.client.Pipeline(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestPipeline(t *testing.T) {
	suite.Run(t, new(SuiteTestPipeline))
}
