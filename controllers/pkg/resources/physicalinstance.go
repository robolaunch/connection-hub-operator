package resources

import (
	"strconv"
	"strings"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getRelayServerSelector(pi connectionhubv1alpha1.PhysicalInstance) map[string]string {
	return map[string]string{
		"relay": pi.Name,
	}
}

var relayPort int = 8080

func GetRelayServerPod(cr *connectionhubv1alpha1.PhysicalInstance) *corev1.Pod {

	apiServerURL := strings.ReplaceAll(cr.Spec.Server, "https://", "")

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetRelayServerPodMetadata().Name,
			Namespace: cr.GetRelayServerPodMetadata().Namespace,
			Labels:    getRelayServerSelector(*cr),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "relay",
					Command: []string{
						"/bin/bash",
						"-c",
						"socat TCP4-LISTEN:" + strconv.Itoa(relayPort) + ",fork,reuseaddr TCP4:" + apiServerURL,
					},
					Image: "robolaunchio/relay:socat-1.7.3.3-2-focal-0.1.0",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: int32(relayPort),
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
		},
	}

	return &pod
}

func GetRelayServerService(cr *connectionhubv1alpha1.PhysicalInstance) *corev1.Service {

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetRelayServerServiceMetadata().Name,
			Namespace: cr.GetRelayServerServiceMetadata().Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     int32(relayPort),
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: int32(relayPort),
					},
				},
			},
			Selector: getRelayServerSelector(*cr),
			Type:     corev1.ServiceTypeNodePort,
		},
	}

	return &svc
}
