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
	yamlStr := ``
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
