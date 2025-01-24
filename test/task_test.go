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
	yamlStr := ``
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

func (s *SuiteTestTask) Test5DeleteTask() {
	err := s.client.Task(s.namespace).Delete(context.TODO(), s.name)
	s.Nil(err)
}

func TestSuiteTestTask(t *testing.T) {
	suite.Run(t, new(SuiteTestTask))
}
