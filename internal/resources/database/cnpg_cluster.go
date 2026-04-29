package database

import (
	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	databasev1alpha1 "platform.io/platform-operator/api/database/v1alpha1"
)

// SyncCNPGCluster maps the CRD PGDatabase to a Cluster from CloudNativePG
func SyncCNPGCluster(pgDb *databasev1alpha1.PGDatabase) *cnpgv1.Cluster {
	return &cnpgv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "postgresql.cnpg.io/v1",
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pgDb.Name,
			Namespace: pgDb.Namespace,
		},
		Spec: cnpgv1.ClusterSpec{
			Instances: int(pgDb.Spec.Instances),
			StorageConfiguration: cnpgv1.StorageConfiguration{
				Size: pgDb.Spec.StorageSize,
			},
			ImageName: "ghcr.io/cloudnative-pg/postgresql:" + pgDb.Spec.Version,
		},
	}
}
