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

package controllers

import (
	"context"
	"github.com/wangwanggogo19xx/devops-golang-test/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("MyStatefulSet controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		MyStsName      = "test-mysts"
		MyStsNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating an  MySts ", func() {
		It("Should be create pods", func() {
			By("By creating a new MySts")
			ctx := context.Background()
			mySts := &v1alpha1.MyStatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      MyStsName,
					Namespace: MyStsNamespace,
				},
				Spec: v1alpha1.MyStatefulSetSpec{
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

			Expect(k8sClient.Create(ctx, mySts)).Should(Succeed())

			createdSts := &v1alpha1.MyStatefulSet{}
			stsName := types.NamespacedName{Namespace: mySts.Namespace, Name: mySts.Name}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, stsName, createdSts)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			updateSts := createdSts
			updateSts.Spec.Replicas = updateSts.Spec.Replicas - 1
			err := k8sClient.Update(ctx, createdSts)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
