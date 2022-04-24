/*
Copyright 2022.

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
	"encoding/json"
	"fmt"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "k8s.io/api/batch/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/go-logr/logr"
	examplecomv1alpha1 "github.com/josefkarasek/distributed-cowboys-operator/api/v1alpha1"
)

const (
	coordinatorNodeID = 0
	configKey         = "config"
	cowboy            = "cowboy"
)

// ShootoutReconciler reconciles a Shootout object
type ShootoutReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Namespace string
	log logr.Logger
}

//+kubebuilder:rbac:groups=example.com,resources=shootouts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=shootouts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=shootouts/finalizers,verbs=update

// Reconcile manages the lifecycle of Shootout
func (r *ShootoutReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log = log.FromContext(ctx)

	instance := &examplecomv1alpha1.Shootout{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if instance.Status.Winner != "" || instance.Status.Error != "" {
		r.log.Info("Shootout is over")
		// this shootout is over
		return ctrl.Result{}, nil
	}

	cowboys, err := unmarshalAndValidateRawInput(instance)
	if err != nil {
		r.log.Error(err, "invalid .spec.cowboys input")
		// non-recoverable
		instance.Status.Error = err.Error()
		return ctrl.Result{}, nil
	}

	if len(cowboys) == 1 {
		// The Invincible One
		r.log.Info("Only one cowboys in this shootout")
		instance.Status.Winner = cowboys[0].Name
		return ctrl.Result{}, nil
	}

	aJob := &batchv1.Job{}
	if err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-0", Namespace: instance.Namespace}, aJob); err != nil {
		if kerrors.IsNotFound(err) {
			r.log.Info("Beginning new shootout")
			err = r.newShootout(ctx, instance, cowboys)
		}
		return ctrl.Result{}, err
	} else {
		// collect results
		r.log.Info("Shootout in progress/over")
	}

	return ctrl.Result{}, nil
}

func unmarshalAndValidateRawInput(instance *examplecomv1alpha1.Shootout) ([]examplecomv1alpha1.Cowboy, error) {
	var cowboys []examplecomv1alpha1.Cowboy
	if err := json.Unmarshal([]byte(instance.Spec.Cowboys), &cowboys); err != nil {
		return nil, err
	}
	if len(cowboys) == 0 {
		return nil, fmt.Errorf("no cowboys taking part in this shootout")
	}
	return cowboys, nil
}

func (r *ShootoutReconciler) newShootout(ctx context.Context, instance *examplecomv1alpha1.Shootout, cowboys []examplecomv1alpha1.Cowboy) error {
	// create config map
	if err := r.Create(ctx, desiredConfigMap(instance)); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			r.log.Error(err, "failed to create config map")
			return err
		}
	}

	var errors []error
	for i, c := range cowboys {
		r.log.Info("Creating cowboy instance", "name", c.Name, "health", c.Health, "damage", c.Damage)
		selectorValue := fmt.Sprintf("%s-%d", instance.Name, i)
		neighbor := fmt.Sprintf("%s-%d", instance.Name, (i+1)%len(cowboys))

		// create service
		if err := r.Create(ctx, desiredService(instance, selectorValue)); err != nil {
			if !kerrors.IsAlreadyExists(err) {
				r.log.Error(err, "failed to create service")
				errors = append(errors, err)
			}
		}
		// create job
		if err := r.Create(ctx, desiredJob(instance, selectorValue, c.Name, neighbor, i == coordinatorNodeID)); err != nil {
			if !kerrors.IsAlreadyExists(err) {
				r.log.Error(err, "failed to create job")
				errors = append(errors, err)
			}
		}
	}

	return utilerrors.NewAggregate(errors)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShootoutReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1alpha1.Shootout{}).
		Complete(r)
}
