package types

type Metadata struct {
	Annotations       map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name              string            `json:"name,omitempty" yaml:"name,omitempty"`
	GenerateName      string            `json:"generateName,omitempty" yaml:"generateName,omitempty"`
	Namespace         string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
}

type TektonResource struct {
	ApiVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   Metadata               `json:"metadata" yaml:"metadata"`
	Spec       map[string]interface{} `json:"spec" yaml:"spec"`
	Status     map[string]interface{} `json:"status" yaml:"status,omitempty"`
}
