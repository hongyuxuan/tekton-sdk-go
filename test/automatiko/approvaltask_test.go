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

type SuiteTestApprovalTask struct {
	suite.Suite
	client    *tekton.Client
	name      string
	namespace string
}

func (s *SuiteTestApprovalTask) SetupSuite() {
	s.client = tekton.NewClient(
		option.WithKubeconfig("../kubeconfig"),
		option.WithSecretPrefix("default-token"),
		// option.WithDebug(true),
	)
	s.name = "test-pipeline-wjx7h-approval"
	s.namespace = "default"
}

func (s *SuiteTestApprovalTask) Test1ListApprovalTask() {
	res, err := s.client.ApprovalTask(s.namespace).List(context.TODO(), metav1.ListOptions{
		Limit: 3,
	})
	s.Nil(err)
	if s.NotNil(res) {
		b, _ := json.Marshal(res)
		fmt.Println(string(b))
	}
}

func (s *SuiteTestApprovalTask) Test2GetApprovalTask() {
	res, err := s.client.ApprovalTask(s.namespace).Get(context.TODO(), s.name)
	s.Nil(err)
	if s.NotNil(res) {
		fmt.Println(res.ToJsonStringPretty())
	}
}

func (s *SuiteTestApprovalTask) Test3GetYamlApprovalTask() {
	res, err := s.client.ApprovalTask(s.namespace).GetYaml(context.TODO(), s.name)
	s.Nil(err)
	if s.NotEmpty(res) {
		fmt.Println(res)
	}
}

func (s *SuiteTestApprovalTask) Test4DeleteApprovalTask() {
	err := s.client.ApprovalTask(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestApprovalTask(t *testing.T) {
	suite.Run(t, new(SuiteTestApprovalTask))
}
