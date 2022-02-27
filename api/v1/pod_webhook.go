/*
Copyright 2022 SiJie.
*/

package v1

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodSidecarMutate struct {
	Client  client.Client
	decoder *admission.Decoder
}

func NewPodSidecarMutate(c client.Client) admission.Handler {
	return &PodSidecarMutate{
		Client: c,
	}
}

func (v *PodSidecarMutate) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := v.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// 添加 sidecar
	sidecar := corev1.Container{
		Name:  "nginx",
		Image: "nginx:1.6",
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 80,
			},
		},
		ImagePullPolicy: corev1.PullIfNotPresent,
	}

	pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (v *PodSidecarMutate) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
