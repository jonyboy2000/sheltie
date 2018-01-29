package structure

import kv1 "k8s.io/api/core/v1"

//============================================
// Config
//============================================

// Config Feilds => Exported
type Config struct {
	K8s_ConfigFilePath string   `json:"k8s_ConfigFilePath"`
	K8s_PodDirectory   string   `json:"k8s_PodDirectory"`
	K8s_NameSpace      string   `json:"k8s_NameSpace"`
	K8s_NodeName       string   `json:"k8s_NodeName"`
	K8s_Priority       int32    `json:"k8s_Priority"`
	K8s_Labels         []string `json:"k8s_Labels"`
	K8s_Annotations    []string `json:"k8s_Annotations"`

	K8s_PodName_Directory string `json:"k8s_PodName_Directory"`
	K8s_PodName_Download  string `json:"k8s_PodName_Download"`

	Rabbit_User         string `json:"rabbit_User"`
	Rabbit_Pass         string `json:"rabbit_Pass"`
	Rabbit_Host         string `json:"rabbit_Host"`
	Rabbit_Port         string `json:"rabbit_Port"`
	Rabbit_CommandQueue string `json:"rabbit_CommandQueue"`
	Rabbit_ResultQueue  string `json:"rabbit_ResultQueue"`
}

//============================================
// K8s Job
//============================================

// YAML file fields => Exported
type YAML struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string
	} `yaml:"metadata"`
	Spec struct {
		Template struct {
			Metadata struct {
				Name string
			} `yaml:"metadata"`
			Spec struct {
				ImagePullSecrets []struct {
					Name string
				} `yaml:"imagePullSecrets"`
				Containers []struct {
					Name         string   `yaml:"name"`
					Image        string   `yaml:"image"`
					Command      []string `yaml:"command"`
					Args         []string `yaml:"args"`
					VolumeMounts []struct {
						MountPath string `yaml:"mountPath"`
						Name      string `yaml:"name"`
					} `yaml:"volumeMounts"`
				} `yaml:"containers"`
				Volumes []struct {
					HostPath struct {
						Path string           `yaml:"path"`
						Type kv1.HostPathType `yaml:"type"`
						Name string           `yaml:"name"`
					} `yaml:"hostPath"`
				} `yaml:"volumes"`
				HostNetwork   bool              `yaml:"hostNetwork"`
				RestartPolicy kv1.RestartPolicy `yaml:"restartPolicy"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}
