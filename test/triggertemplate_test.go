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

type SuiteTestTriggerTemplate struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestTriggerTemplate) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "testtriggertemplate"
	s.namespace = "default"
}

func (s *SuiteTestTriggerTemplate) Test1CreateTriggerTemplate() {
	yamlStr := `apiVersion: triggers.tekton.dev/v1beta1
kind: TriggerTemplate
metadata:
  annotations:
    fiops/author: hongyuxuan
  labels:
    app: testtriggertemplate
  creationTimestamp: "2024-03-07T01:31:16Z"
  name: testtriggertemplate
  namespace: default
spec:
  params:
  - default: master
    description: The git revision
    name: gitrevision
  - description: The git repository url
    name: gitrepositoryurl
  - default: tekton-pipeline
    description: The namespace to create the resources
    name: namespace
  - description: The application needs to be triggered
    name: application
  resourcetemplates:
  - apiVersion: tekton.dev/v1beta1
    kind: PipelineRun
    metadata:
      generateName: $(tt.params.application)-pipeline-run-
      namespace: $(tt.params.namespace)
    spec:
      params:
      - name: revision
        value: $(tt.params.gitrevision)
      - name: repo-url
        value: $(tt.params.gitrepositoryurl)
      pipelineRef:
        name: $(tt.params.application)-pipeline
      serviceAccountName: default
      workspaces:
      - name: shared-workspace
        volumeClaimTemplate:
          spec:
            accessModes:
            - ReadWriteOnce
            resources:
              requests:
                storage: 1Gi
            storageClassName: nfs-client
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
	err := s.client.TriggerTemplate(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestTriggerTemplate) Test2ListTriggerTemplate() {
	res, err := s.client.TriggerTemplate(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=testtriggertemplate",
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

func (s *SuiteTestTriggerTemplate) Test3GetTriggerTemplate() {
	res, err := s.client.TriggerTemplate(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestTriggerTemplate) Test4GetYamlTriggerTemplate() {
	res, err := s.client.TriggerTemplate(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestTriggerTemplate) Test5DeleteTriggerTemplate() {
	err := s.client.TriggerTemplate(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestTriggerTemplate(t *testing.T) {
	suite.Run(t, new(SuiteTestTriggerTemplate))
}
