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

type SuiteTestPipelineRun struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestPipelineRun) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "testpipelinerun"
	s.namespace = "default"
}

func (s *SuiteTestPipelineRun) Test1CreatePipelineRun() {
	yamlStr := `apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  annotations:
    fiops/author: hongyuxuan
  creationTimestamp: "2024-09-04T02:08:26Z"
  labels:
    app: testpipelinerun
    dashboard.tekton.dev/rerunOf: lizardrestic-server-pipeline-run-r-v9lc7
    tekton.dev/pipeline: lizardrestic-server-pipeline
    triggers.tekton.dev/eventlistener: fiops-pipeline-eventlistener
    triggers.tekton.dev/trigger: fiops-match
    triggers.tekton.dev/triggers-eventid: 80284306-1308-4adc-9574-7ddb706ca4d4
  name: testpipelinerun
  generateName: lizardrestic-server-pipeline-run-r-
  namespace: default
spec:
  params:
  - name: revision
    value: release-v1.0.0
  - name: repo-url
    value: git@gitlab.cicconline.com:xficc/devops/lizardrestic.git
  pipelineRef:
    name: lizardrestic-server-pipeline
  taskRunTemplate:
    serviceAccountName: default
  timeouts:
    pipeline: 1h0m0s
  workspaces:
  - name: shared-workspace
    volumeClaimTemplate:
      metadata:
        creationTimestamp: null
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
        storageClassName: nfs-client
      status: {}
  - name: dockerhub-auth
    secret:
      secretName: docker-credential
  - name: git-credentials
    secret:
      secretName: git-credentials
  - name: jfrog-auth
    secret:
      secretName: jfrog-credentials
`
	err := s.client.PipelineRun(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestPipelineRun) Test2ListPipelineRun() {
	res, err := s.client.PipelineRun(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=testpipelinerun",
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

func (s *SuiteTestPipelineRun) Test3GetPipelineRun() {
	res, err := s.client.PipelineRun(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipelineRun) Test4GetYamlPipelineRun() {
	res, err := s.client.PipelineRun(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipelineRun) Test5DeletePipelineRun() {
	err := s.client.PipelineRun(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestPipelineRun(t *testing.T) {
	suite.Run(t, new(SuiteTestPipelineRun))
}
