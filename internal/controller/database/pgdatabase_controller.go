/*
Copyright 2026.

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

package database

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	databasev1alpha1 "platform.io/platform-operator/api/database/v1alpha1"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PGDatabaseReconciler reconciles a PGDatabase object
type PGDatabaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.yourdomain.com,resources=pgdatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.yourdomain.com,resources=pgdatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=postgresql.cnpg.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete

func (r *PGDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// 1. Fetch the PGDatabase instance
	var pgDb databasev1alpha1.PGDatabase
	if err := r.Get(ctx, req.NamespacedName, &pgDb); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Define the desired CNPG Cluster object
	desiredCluster := &cnpgv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pgDb.Name,
			Namespace: pgDb.Namespace,
		},
		Spec: cnpgv1.ClusterSpec{
			Instances: int(pgDb.Spec.Instances),
			StorageConfiguration: cnpgv1.StorageConfiguration{
				Size: pgDb.Spec.StorageSize,
			},
			// Simplification: Using the version as the image tag
			ImageName: "ghcr.io/cloudnative-pg/postgresql:" + pgDb.Spec.Version,
		},
	}

	// 3. Set PGDatabase as the owner of the Cluster
	// This ensures that if PGDatabase is deleted, the Cluster is also deleted
	if err := ctrl.SetControllerReference(&pgDb, desiredCluster, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// 4. Logic to Create or Update the Cluster (Simplified Server-Side Apply style)
	var existingCluster cnpgv1.Cluster
	err := r.Get(ctx, client.ObjectKey{Name: desiredCluster.Name, Namespace: desiredCluster.Namespace}, &existingCluster)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			l.Info("Creating a new CNPG Cluster")
			return ctrl.Result{}, r.Create(ctx, desiredCluster)
		}
		return ctrl.Result{}, err
	}

	// Logic to update existingCluster would go here if needed...

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PGDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.PGDatabase{}).
		Owns(&cnpgv1.Cluster{}). // Watch Clusters owned by PGDatabase
		Complete(r)
}
