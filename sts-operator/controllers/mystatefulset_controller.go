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
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	stsv1alpha1 "github.com/wangwanggogo19xx/devops-golang-test/api/v1alpha1"
)

// MyStatefulSetReconciler reconciles a MyStatefulSet object
type MyStatefulSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sts.example.com,resources=mystatefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sts.example.com,resources=mystatefulsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sts.example.com,resources=mystatefulsets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyStatefulSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MyStatefulSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return r.ReconcileMyStatefulSet(ctx, req)

}

// SetupWithManager sets up the controller with the Manager.
func (r *MyStatefulSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stsv1alpha1.MyStatefulSet{}).
		Complete(r)
}

func (r *MyStatefulSetReconciler) ReconcileMyStatefulSet(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mySts := &stsv1alpha1.MyStatefulSet{}
	err := r.Get(ctx, req.NamespacedName, mySts)
	if err != nil {
		return ctrl.Result{}, nil
	}
	currentPod := &v1.Pod{}

	var podName string
	for i := 0; i < mySts.Spec.Replicas; i++ {

		podName = fmt.Sprintf("%v-%d", mySts.Name, i)
		err := r.Client.Get(ctx, types.NamespacedName{Namespace: mySts.Namespace, Name: podName}, currentPod)
		if err != nil && errors.IsNotFound(err) {
			log.Log.Info(fmt.Sprintf("creating %v", podName))
			expectedPod, _ := r.getExpectedStsPod(mySts, podName)
			err := r.Create(ctx, expectedPod)
			return ctrl.Result{RequeueAfter: time.Duration(15) * time.Second}, err
		}

		if currentPod.Status.Phase != v1.PodRunning {
			log.Log.Info(fmt.Sprintf("waiting pod(%v) Runnnig", podName))
			return ctrl.Result{RequeueAfter: time.Duration(15) * time.Second}, err
		} else {
			//_ = append(mySts.Status.PodStatus, currentPod.Status)
			//if mySts.Status.PodStatus == nil {
			//	mySts.Status.PodStatus = []v1.PodStatus{}
			//}
			//podStatus := mySts.Status.PodStatus[i]
			if len(mySts.Status.PodStatus) <= i {
				mySts.Status.PodStatus = append(mySts.Status.PodStatus, currentPod.Status)
				r.Status().Update(ctx, mySts)
			}

			expectedPod, _ := r.getExpectedStsPod(mySts, podName)

			if r.needUpdated(currentPod, expectedPod) {
				log.Log.Info(fmt.Sprintf("update  pod(%v) ", podName))
				err := r.Update(ctx, expectedPod)
				return ctrl.Result{Requeue: true}, err
			}
		}
	}
	myStateFulSetPodLabel := r.getExpectedStsLabels(mySts)
	selector := labels.SelectorFromSet(myStateFulSetPodLabel)

	option := &client.ListOptions{LabelSelector: selector, Namespace: mySts.Namespace}
	podList := &v1.PodList{}
	err = r.List(ctx, podList, option)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Duration(15) * time.Second}, err
	}

	for i := len(podList.Items) - 1; i >= mySts.Spec.Replicas; i++ {
		deletePodName := fmt.Sprintf("%v-%d", mySts.Name, i)
		deletePod := &v1.Pod{
			ObjectMeta: v12.ObjectMeta{
				Name:      deletePodName,
				Namespace: mySts.Namespace,
			},
		}
		err := r.Delete(ctx, deletePod)
		if err != nil {
			log.Log.Info(err.Error())
		}
		return ctrl.Result{RequeueAfter: time.Duration(15) * time.Second}, err
	}

	mySts.Status.Phase = "Running"
	r.Status().Update(ctx, mySts)
	return ctrl.Result{}, err
}

func (r *MyStatefulSetReconciler) getExpectedStsPod(mySts *stsv1alpha1.MyStatefulSet, podName string) (*v1.Pod, error) {
	myStateFulSetPodLabel := r.getExpectedStsLabels(mySts)

	expectedPod := &v1.Pod{
		ObjectMeta: v12.ObjectMeta{
			Name:      podName,
			Namespace: mySts.Namespace,
			Labels:    myStateFulSetPodLabel,
			OwnerReferences: []v12.OwnerReference{
				{
					APIVersion:         mySts.APIVersion,
					Kind:               mySts.Kind,
					Name:               mySts.Name,
					UID:                mySts.UID,
					BlockOwnerDeletion: pointer.Bool(true),
					Controller:         pointer.Bool(true),
				},
			},
		},
		Spec: mySts.Spec.Template.Spec,
	}

	return expectedPod, nil
}

func (r *MyStatefulSetReconciler) getExpectedStsLabels(mysts *stsv1alpha1.MyStatefulSet) map[string]string {
	myStateFulSetPodLabel := mysts.Labels
	if myStateFulSetPodLabel == nil {
		myStateFulSetPodLabel = make(map[string]string)
	}
	myStateFulSetPodLabel[stsv1alpha1.LABEL_MYSTS_CONTROLLER_KEY] = stsv1alpha1.LABEL_MYSTS_CONTROLLER_VALUE

	return myStateFulSetPodLabel
}

func (r *MyStatefulSetReconciler) needUpdated(currentPod *v1.Pod, expectedPod *v1.Pod) bool {
	log.Log.Info(fmt.Sprintf("check pod %v Check for updates", currentPod.Name))
	expectedPodLabels := expectedPod.GetLabels()
	currentPodLabels := currentPod.GetLabels()

	if expectedPodLabels == nil || currentPodLabels == nil {
		return true
	}
	for key, value := range expectedPodLabels {
		if currentPodLabels[key] != value {
			return true
		}
	}
	return false
}
