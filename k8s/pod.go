package k8s

import (
	st "../structure"
	util "../utility"
	yaml "gopkg.in/yaml.v2"

	kv1 "k8s.io/api/core/v1"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k "k8s.io/client-go/kubernetes"
	kclient "k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// CreatePOD => Exported
func CreatePOD(
	kubeConfigPath string,
	podFilePath string, // type = YAML file
	namespace string,
	nodeName string,
	priority int32,
	labels map[string]string,
	annotations map[string]string,
	commands [][]string, // [container_index][command_index]
	args [][]string) bool { // [container_index][arg_index]
	ret := false
	var err error

	//============================================
	// Parse POD File
	//============================================
	var pod st.YAML
	content := util.ReadFile(podFilePath)
	if content == nil {
		util.Log("(K8s=>CreatePOD) Error in read pod file", nil)
		return ret
	}

	err = yaml.Unmarshal(content, &pod)
	if err != nil {
		util.Log("(K8s=>CreatePOD) Error in unmarshalling pod file", err)
		return ret
	}

	//============================================
	// Create Kubernetes Config
	//============================================

	// Use the current context in kubeconfig
	config, err := kclient.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		// panic(err.Error())
		util.Log("(K8s=>CreatePOD) Error in building k8s config", err)
		return ret
	}

	// Create the clientset
	clientset, err := k.NewForConfig(config)
	if err != nil {
		// panic(err.Error())
		util.Log("(K8s=>CreatePOD) Error in building k8s clientset", err)
		return ret
	}

	//===================================================
	// Find and Delete POD
	//===================================================
	podName := pod.Metadata.Name
	deletionGraceTime := int64(0)

	_, err = clientset.CoreV1().Pods(namespace).Get(podName, kmetav1.GetOptions{})
	if err == nil {
		err = clientset.CoreV1().Pods(namespace).Delete(podName, &kmetav1.DeleteOptions{GracePeriodSeconds: &deletionGraceTime})
		if err != nil {
			//panic(err.Error())
			util.Log("(K8s=>CreatePOD) Error deleting pod", err)
			return ret
		}
	} else if (err != nil) && (!kerror.IsNotFound(err)) {
		//panic(err.Error())
		util.Log("(K8s=>CreatePOD) Error getting pod", err)
		return ret
	}

	//===================================================
	// Generate POD Data
	//===================================================

	// ImagePullSecrets
	//=====================
	_imagePullSecrets := pod.Spec.Template.Spec.ImagePullSecrets
	imagePullSecrets := []kv1.LocalObjectReference{}

	for i := 0; i < len(_imagePullSecrets); i++ {
		imagePullSecrets = append(imagePullSecrets, kv1.LocalObjectReference{Name: _imagePullSecrets[i].Name})
	}

	// Volumes
	//=====================
	_volumes := pod.Spec.Template.Spec.Volumes
	volumes := []kv1.Volume{}

	for i := 0; i < len(_volumes); i++ {
		_type := _volumes[i].HostPath.Type
		_hostPath := kv1.HostPathVolumeSource{
			Path: _volumes[i].HostPath.Path,
			Type: &_type,
		}

		volumes = append(volumes, kv1.Volume{
			Name:         _volumes[i].HostPath.Name,
			VolumeSource: kv1.VolumeSource{HostPath: &_hostPath}})
	}

	// Containers
	//=====================
	_containers := pod.Spec.Template.Spec.Containers
	containers := []kv1.Container{}

	for i := 0; i < len(_containers); i++ {

		_volumeMounts := _containers[i].VolumeMounts
		volumeMounts := []kv1.VolumeMount{}
		for i := 0; i < len(_volumeMounts); i++ {
			volumeMounts = append(volumeMounts, kv1.VolumeMount{
				Name:      _volumeMounts[i].Name,
				MountPath: _volumeMounts[i].MountPath,
			})
		}

		_command := commands[i]
		for k := 0; k < len(_containers[i].Command); k++ {
			_command = append(_command, _containers[i].Command[k])
		}

		_args := args[i]
		for k := 0; k < len(_containers[i].Args); k++ {
			_args = append(_args, _containers[i].Args[k])
		}

		containers = append(containers, kv1.Container{
			Name:         _containers[i].Name,
			Image:        _containers[i].Image,
			Command:      _command,
			Args:         _args,
			VolumeMounts: volumeMounts,
		})
	}

	// Other
	//=====================
	restartPolicy := pod.Spec.Template.Spec.RestartPolicy
	hostNetwork := pod.Spec.Template.Spec.HostNetwork
	APIVersion := pod.APIVersion
	kind := pod.Kind

	//===================================================
	// Create POD
	//===================================================
	p := &kv1.Pod{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:        podName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: kv1.PodSpec{
			Affinity:         &kv1.Affinity{},
			RestartPolicy:    restartPolicy,
			HostNetwork:      hostNetwork,
			NodeName:         nodeName,
			ImagePullSecrets: imagePullSecrets,
			Priority:         &priority,
			Volumes:          volumes,
			Containers:       containers,
		},
	}
	p.APIVersion = APIVersion
	p.Kind = kind

	_, err = clientset.CoreV1().Pods(namespace).Create(p)
	if err != nil {
		//panic(err.Error())
		util.Log("(K8s=>CreatePOD) Error creating pod", err)
	} else {
		ret = true
	}

	//===================================================
	// Return
	//===================================================
	return ret
}
