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
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mystatefulsetlog = logf.Log.WithName("mystatefulset-resource")

func (r *MyStatefulSet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-sts-example-com-example-com-v1alpha1-mystatefulset,mutating=true,failurePolicy=fail,sideEffects=None,groups=sts.example.com.example.com,resources=mystatefulsets,verbs=create;update,versions=v1alpha1,name=mmystatefulset.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MyStatefulSet{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MyStatefulSet) Default() {
	mystatefulsetlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Labels == nil {
		r.Labels = make(map[string]string)
	}
	val, ok := r.Labels[LABEL_MYSTS_CONTROLLER_KEY]
	if !ok || val != LABEL_MYSTS_CONTROLLER_VALUE {
		r.Labels[LABEL_MYSTS_CONTROLLER_KEY] = LABEL_MYSTS_CONTROLLER_VALUE
	}

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-sts-example-com-example-com-v1alpha1-mystatefulset,mutating=false,failurePolicy=fail,sideEffects=None,groups=sts.example.com.example.com,resources=mystatefulsets,verbs=create;update,versions=v1alpha1,name=vmystatefulset.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MyStatefulSet{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MyStatefulSet) ValidateCreate() error {
	mystatefulsetlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.

	return r.checkLabels()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MyStatefulSet) ValidateUpdate(old runtime.Object) error {
	mystatefulsetlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.checkLabels()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MyStatefulSet) ValidateDelete() error {
	mystatefulsetlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return r.checkLabels()
}

func (r *MyStatefulSet) checkLabels() error {
	if r.Labels == nil {
		return errors.New("labels can not be nil")
	}
	return nil
}
