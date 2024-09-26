/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("MyStatefulSet Webhook", func() {

	Context("When creating MyStatefulSet under Defaulting Webhook", func() {
		It("Should fill default labels", func() {
			const (
				MyStsName      = "test-mysts"
				MyStsNamespace = "default"
			)
			// TODO(user): Add your logic here

			mySts := &MyStatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      MyStsName,
					Namespace: MyStsNamespace,
				},
				Spec: MyStatefulSetSpec{
					Replicas: 3,
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: MyStsName,
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "nginx",
									Image: "nginx",
								},
							},
						},
					},
				},
			}

			mySts.ValidateCreate()
			Expect(mySts.ValidateCreate()).To(HaveOccurred())
			mySts.ValidateDelete()

			mySts.Default()
			Expect(mySts.Labels[LABEL_MYSTS_CONTROLLER_KEY]).Should(Equal(LABEL_MYSTS_CONTROLLER_VALUE))

		})
	})

})
