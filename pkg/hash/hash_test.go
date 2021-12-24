/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package hash

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetTemplateHash(t *testing.T) {
	h := NewTappHash()

	template := createPodTemplate()
	expectedTemplateHash := generateTemplateHash(&template)
	h.SetTemplateHash(&template)
	realHash := h.GetTemplateHash(template.Labels)
	if expectedTemplateHash != realHash {
		t.Errorf("Failed to set template h")
	}
}

func TestSetUniqHash(t *testing.T) {
	h := NewTappHash()

	template := createPodTemplate()
	expectedUniqHash := generateUniqHash(template)
	h.SetUniqHash(&template)
	realHash := h.GetUniqHash(template.Labels)
	if expectedUniqHash != realHash {
		t.Errorf("Failed to set uniq h")
	}
}

func TestSetSpecHash(t *testing.T) {
	h := NewTappHash()

	template := createPodTemplate()
	expectedUniqHash := generateSpecHash(template)
	h.SetSpecHash(&template)
	realHash := h.GetSpecHash(template.Labels)
	if expectedUniqHash != realHash {
		t.Errorf("Failed to set uniq h")
	}
}

func createPodTemplate() corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{"test": "hello"},
			Annotations: map[string]string{},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyOnFailure,
			DNSPolicy:     corev1.DNSClusterFirst,
			Containers:    []corev1.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
		},
	}
}
