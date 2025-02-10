package automatiko

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApprovalTask struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds the desired state of the ApprovalTask from the client
	// +optional
	Spec ApprovalTaskSpec `json:"spec"`
	// +optional
	Status ApprovalTaskStatus `json:"status,omitempty"`
}

type ApprovalTaskSpec struct {
	Approvers   []string `json:"approvers"`
	Description string   `json:"description"`
	Groups      []string `json:"groups"`
	Pipeline    string   `json:"pipeline"`
	Strategy    string   `json:"strategy"`
}

type ApprovalTaskStatus struct {
	ApprovalUrl string              `json:"approvalUrl"`
	Message     string              `json:"message"`
	Results     ApprovalTaskResults `json:"results,omitempty"`
	Status      string              `json:"status"`
}

type ApprovalTaskResults struct {
	ApprovedBy string `json:"approvedBy,omitempty"`
	RejectedBy string `json:"rejectedBy,omitempty"`
	Comment    string `json:"comment"`
	Decision   string `json:"decision"`
}

func (t *ApprovalTask) ToJsonString() string {
	b, _ := json.Marshal(t)
	return string(b)
}

func (t *ApprovalTask) ToJsonStringPretty() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}
