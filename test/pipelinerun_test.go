package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	tekton "github.com/hongyuxuan/tekton-sdk-go"
	"github.com/hongyuxuan/tekton-sdk-go/core/option"
	"github.com/hongyuxuan/tekton-sdk-go/types"
	"github.com/stretchr/testify/suite"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
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
		option.WithKubeconfig("/root/.kube/config"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "lizardcd-ui-m5nkw-r-v5spb"
	s.namespace = "tektoncd-default"
}

func (s *SuiteTestPipelineRun) Test1CreatePipelineRun() {
	yamlStr := ``
	err := s.client.PipelineRun(s.namespace).Create(context.TODO(), yamlStr)
	s.Nil(err)
}

func (s *SuiteTestPipelineRun) Test2ListPipelineRun() {
	res, _, err := s.client.PipelineRun(s.namespace).List(context.TODO(), metav1.ListOptions{
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

func (s *SuiteTestPipelineRun) Test3GetPipelineRun() {
	res, err := s.client.PipelineRun(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(b))
	}
}

func (s *SuiteTestPipelineRun) Test4CancelPipelineRun() {
	res, err := s.client.PipelineRun(s.namespace).Patch(context.TODO(), s.name, []types.PatchOptions{
		{
			Op:    "replace",
			Path:  "/spec/status",
			Value: "Cancelled",
		},
	})
	s.Nil(err)
	if s.NotNil(res) {
		s.Equal(tektonv1.PipelineRunSpecStatus("Cancelled"), res.Spec.Status)
	}
}

func (s *SuiteTestPipelineRun) Test5GetYamlPipelineRun() {
	res, err := s.client.PipelineRun(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestPipelineRun) Test6DeletePipelineRun() {
	err := s.client.PipelineRun(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestPipelineRun(t *testing.T) {
	suite.Run(t, new(SuiteTestPipelineRun))
}
