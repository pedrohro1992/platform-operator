# Platform Operator

Platform Operator is a Kubernetes operator designed to abstract and simplify the provisioning of platform infrastructure components such as Databases and Secret Management. It acts as an Internal Developer Platform (IDP) enabler by providing high-level Custom Resource Definitions (CRDs) that wrap underlying complex operators.

## Supported Custom Resources

The operator currently supports the following Custom Resources:

### 1. PGDatabase (`database.platform.io/v1alpha1`)
Provisions a high-availability PostgreSQL database leveraging the [CloudNativePG](https://cloudnative-pg.io/) (CNPG) operator.

**Spec Fields:**
- `instances` (int32): Number of Postgres instances (Default: `3`).
- `storageSize` (string): Size of the storage (e.g., `"1Gi"`).
- `version` (string): PostgreSQL major version (Default: `"15"`).

**Under the hood:** The controller watches `PGDatabase` resources and reconciles them into CNPG `Cluster` resources using Server-Side Apply, ensuring correct ownership and lifecycle management.

### 2. VaultConnection (`security.platform.io/v1alpha1`)
Configures a cluster-wide connection to HashiCorp Vault using the [External Secrets Operator](https://external-secrets.io/) (ESO).
*Note: This is a Cluster-scoped resource.*

**Spec Fields:**
- `vaultUrl` (string): The URL to the Vault server.
- `authPath` (string): The Vault authentication path.
- `mountPath` (string): The Vault secrets mount path.
- `vaultRole` (string): The Vault role to authenticate against.

**Under the hood:** The controller watches `VaultConnection` resources and translates them into an ESO `ClusterSecretStore` object to facilitate cluster-wide secret synchronization.

## Prerequisites & Dependencies

Before using the Platform Operator, ensure the following operators are installed in your cluster:
- **CloudNativePG** (To reconcile CNPG `Cluster` objects)
- **External Secrets Operator** (To reconcile `ClusterSecretStore` objects)

## Getting Started

### Prerequisites
- Go version v1.22.0+
- Docker version 17.03+
- kubectl version v1.11.3+
- Access to a Kubernetes v1.11.3+ cluster

### Local Development
To test the controller locally against your current kubeconfig context without deploying the manager pods:
```sh
make install
make run
```

### Deployment to a Cluster
**1. Build and push the operator image:**
```sh
make docker-build docker-push IMG=<your-registry>/platform-operator:tag
```

**2. Install the CRDs into the cluster:**
```sh
make install
```

**3. Deploy the operator to the cluster:**
```sh
make deploy IMG=<your-registry>/platform-operator:tag
```

### Cleanup
**Undeploy the controller and delete CRDs:**
```sh
make undeploy
make uninstall
```
