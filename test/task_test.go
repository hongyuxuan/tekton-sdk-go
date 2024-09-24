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
	s.name = "testtask"
	s.namespace = "default"
}

func (s *SuiteTestTask) Test1CreateTask() {
	yamlStr := `apiVersion: tekton.dev/v1
kind: Task
metadata:
  annotations:
    fiops/author: fanpengfei
    tekton.dev/categories: CLI
    tekton.dev/displayName: Tekton CLI
    tekton.dev/pipelines.minVersion: 0.17.0
    tekton.dev/platforms: linux/amd64,linux/s390x,linux/ppc64le
    tekton.dev/tags: cli
  labels:
    app: testtask
    app.kubernetes.io/version: "0.4"
  name: testtask
  namespace: default
spec:
  description: This task performs operations on Tekton resources using tkn
  params:
  - default: repo.cicc.com.cn/public-docker-virtual/tekton-releases/dogfooding/tkn@sha256:d17fec04f655551464a47dd59553c9b44cf660cc72dbcdbd52c0b8e8668c0579
    description: tkn CLI container image to run this task
    name: TKN_IMAGE
    type: string
  - default: tkn $@
    description: tkn CLI script to execute
    name: SCRIPT
    type: string
  - default:
    - --help
    description: tkn CLI arguments to run
    name: ARGS
    type: array
  steps:
  - args:
    - $(params.ARGS)
    computeResources: {}
    env:
    - name: HOME
      value: /tekton/home
    image: $(params.TKN_IMAGE)
    name: tkn
    script: |
      if [ "$(workspaces.kubeconfig.bound)" = "true" ] && [ -e $(workspaces.kubeconfig.path)/kubeconfig ]; then
        export KUBECONFIG="$(workspaces.kubeconfig.path)"/kubeconfig
      fi

      eval "$(params.SCRIPT)"
    securityContext:
      runAsNonRoot: true
      runAsUser: 65532
  workspaces:
  - description: An optional workspace that allows you to provide a .kube/config file
      for tkn to access the cluster. The file should be placed at the root of the
      Workspace with name kubeconfig.
    name: kubeconfig
    optional: true
`
	err := s.client.Task("default").Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestTask) Test2ListTask() {
	res, err := s.client.Task(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=testtask",
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

func (s *SuiteTestTask) Test5DeleteTask() {
	err := s.client.Task(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestTask(t *testing.T) {
	suite.Run(t, new(SuiteTestTask))
}
