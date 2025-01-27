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

type SuiteTestTaskRun struct {
	suite.Suite
	client    *tekton.Client
	namespace string
	name      string
}

func (s *SuiteTestTaskRun) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("./kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.namespace = "default"
}

func (s *SuiteTestTaskRun) Test1ListTaskRun() {
	res, err := s.client.TaskRun(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "tekton.dev/pipelineRun=maven-pipeline-run-8g47s-r-m5cdd",
	})
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res.ToJsonString())
		s.name = res[0].Name
	}
}

func (s *SuiteTestTaskRun) Test2GetTaskRun() {
	_, err := s.client.TaskRun(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
}

func (s *SuiteTestTaskRun) Test2GetTaskRunStatus() {
	manifest, err := s.client.TaskRun(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	fmt.Println(manifest)
}

func TestSuiteTestTaskRun(t *testing.T) {
	suite.Run(t, new(SuiteTestTaskRun))
}
