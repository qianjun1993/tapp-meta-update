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
	"fmt"
	"hash/fnv"

	corev1 "k8s.io/api/core/v1"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

const (
	// TemplateHashKey is a key for storing PodTemplateSpec's hash value in labels.
	// It will be used to check whether pod's PodTemplateSpec has changed, if yes,
	// we need recreate or do in-place update for the pod according to the value of UniqHash.
	TemplateHashKey = "tapp_template_hash_key"
	// UniqHashKey is a key for storing hash value of PodTemplateSpec(without container images) in labels.
	// It will will be used to check whether pod's PodTemplateSpec hash changed and only container images
	// changed, if yes, we will do in place update for the pod.
	UniqHashKey = "tapp_uniq_hash_key"
	// SpecHashKey is a key for storing hash value of PodTemplateSpec(with container images) in labels.
	// It will will be used to check whether pod's PodTemplateSpec hash changed and only meta
	// changed, if yes, we will do update for the pod.
	SpecHashKey = "tapp_spec_hash_key"
)

// TappHashInterface is used for generate and verify hash for tapp.
type TappHashInterface interface {
	// SetTemplateHash sets PodTemplateSpec's hash value into template's labels,
	// returns true if needs set and is set, otherwise false
	SetTemplateHash(template *corev1.PodTemplateSpec) bool
	// GetTemplateHash returns PodTemplateSpec's hash value, the values is stored in labels.
	GetTemplateHash(labels map[string]string) string
	// SetUniqHash sets hash value of PodTemplateSpec(without container images) into template's labels,
	// returns true if needs set and is set, otherwise false
	SetUniqHash(template *corev1.PodTemplateSpec) bool
	// GetUniqHash returns hash value of PodTemplateSpec(without container images), the values is stored in labels.
	GetUniqHash(labels map[string]string) string
	// SetSpecHash sets hash value of PodTemplateSpec(with container images) into template's labels,
	// returns true if needs set and is set, otherwise false
	SetSpecHash(template *corev1.PodTemplateSpec) bool
	// GetSpecHash returns hash value of PodTemplateSpec(with container images), the values is stored in labels.
	GetSpecHash(labels map[string]string) string
	// HashLabels returns labels key that stores TemplateHash and UniqHash
	HashLabels() []string
}

func NewTappHash() TappHashInterface {
	return &defaultTappHash{}
}

type defaultTappHash struct{}

func (th *defaultTappHash) SetTemplateHash(template *corev1.PodTemplateSpec) bool {
	expected := generateTemplateHash(template)
	hash := th.GetTemplateHash(template.Labels)
	if hash != expected {
		if template.Labels == nil {
			template.Labels = make(map[string]string)
		}
		template.Labels[TemplateHashKey] = expected
		return true
	} else {
		return false
	}
}

func (th *defaultTappHash) GetTemplateHash(labels map[string]string) string {
	return labels[TemplateHashKey]
}

func (th *defaultTappHash) SetUniqHash(template *corev1.PodTemplateSpec) bool {
	expected := generateUniqHash(*template)
	hash := th.GetUniqHash(template.Labels)
	if hash != expected {
		if template.Labels == nil {
			template.Labels = make(map[string]string)
		}
		template.Labels[UniqHashKey] = expected
		return true
	} else {
		return false
	}
}

func (th *defaultTappHash) GetUniqHash(labels map[string]string) string {
	return labels[UniqHashKey]
}

func (th *defaultTappHash) SetSpecHash(template *corev1.PodTemplateSpec) bool {
	expected := generateSpecHash(*template)
	hash := th.GetSpecHash(template.Labels)
	if hash != expected {
		if template.Labels == nil {
			template.Labels = make(map[string]string)
		}
		template.Labels[SpecHashKey] = expected
		return true
	} else {
		return false
	}
}

func (th *defaultTappHash) GetSpecHash(labels map[string]string) string {
	return labels[SpecHashKey]
}

func (th *defaultTappHash) HashLabels() []string {
	return []string{TemplateHashKey, UniqHashKey, SpecHashKey}
}

func generateHash(template interface{}) uint64 {
	hasher := fnv.New64()
	hashutil.DeepHashObject(hasher, template)
	return hasher.Sum64()
}

func generateTemplateHash(template *corev1.PodTemplateSpec) string {
	meta := template.ObjectMeta.DeepCopy()
	delete(meta.Labels, TemplateHashKey)
	delete(meta.Labels, UniqHashKey)
	return fmt.Sprintf("%d", generateHash(corev1.PodTemplateSpec{
		ObjectMeta: *meta,
		Spec:       template.Spec,
	}))
}

func generateUniqHash(template corev1.PodTemplateSpec) string {
	if template.Spec.InitContainers != nil {
		var newContainers []corev1.Container
		for _, container := range template.Spec.InitContainers {
			container.Image = ""
			newContainers = append(newContainers, container)
		}
		template.Spec.InitContainers = newContainers
	}

	var newContainers []corev1.Container
	for _, container := range template.Spec.Containers {
		container.Image = ""
		newContainers = append(newContainers, container)
	}
	template.Spec.Containers = newContainers

	return fmt.Sprintf("%d", generateHash(template.Spec))
}

func generateSpecHash(template corev1.PodTemplateSpec) string {
	return fmt.Sprintf("%d", generateHash(template.Spec))
}
