package security

import (
	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	esmeta "github.com/external-secrets/external-secrets/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	securityv1alpha1 "platform.io/platform-operator/api/security/v1alpha1"
)

// SuncClusterStore maps the CRD VaultConnection to a ClusterSecretStore from ESO
func SyncClusterSecretStore(vc *securityv1alpha1.VaultConnection) *esv1.ClusterSecretStore {
	return &esv1.ClusterSecretStore{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "external-secrets.io/v1",
			Kind:       "ClusterSecretStore",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: vc.Name,
		},
		Spec: esv1.SecretStoreSpec{
			Provider: &esv1.SecretStoreProvider{
				Vault: &esv1.VaultProvider{
					Server:  vc.Spec.VaultUrl,
					Path:    ptr.To(vc.Spec.MountPath),
					Version: esv1.VaultKVStoreV2,
					Auth: &esv1.VaultAuth{
						Kubernetes: &esv1.VaultKubernetesAuth{
							Path: vc.Spec.AuthPath,
							Role: vc.Spec.VaultRole,
							ServiceAccountRef: &esmeta.ServiceAccountSelector{
								Name: "eso-role",
							},
						},
					},
				},
			},
		},
	}
}
