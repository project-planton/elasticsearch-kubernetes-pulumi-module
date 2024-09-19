## Usage

### Sample YAML Configuration

Create a YAML file (`elasticsearch-cluster.yaml`) with the desired configuration:

```yaml
apiVersion: code2cloud.planton.cloud/v1
kind: ElasticsearchKubernetes
metadata:
  id: my-elasticsearch-cluster
spec:
  kubernetes_cluster_credential_id: your-cluster-credential-id
  elasticsearch_container:
    replicas: 3
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "1"
        memory: "2Gi"
    is_persistence_enabled: true
    disk_size: "20Gi"
  kibana_container:
    is_enabled: true
    replicas: 1
    resources:
      requests:
        cpu: "200m"
        memory: "512Mi"
      limits:
        cpu: "500m"
        memory: "1Gi"
  ingress:
    is_enabled: true
    endpoint_domain_name: "example.com"
```

- **`kubernetes_cluster_credential_id`**: The ID of your Kubernetes cluster credentials.
- **`elasticsearch_container`**:
    - **`replicas`**: Number of Elasticsearch pods to deploy.
    - **`resources`**: CPU and memory resources for the Elasticsearch container.
    - **`is_persistence_enabled`**: Set to `true` to enable data persistence.
    - **`disk_size`**: Size of the persistent volume (e.g., "20Gi").
- **`kibana_container`**:
    - **`is_enabled`**: Set to `true` to deploy Kibana.
    - **`replicas`**: Number of Kibana pods to deploy.
    - **`resources`**: CPU and memory resources for the Kibana container.
- **`ingress`**:
    - **`is_enabled`**: Set to `true` to enable ingress configurations.
    - **`endpoint_domain_name`**: The domain name for external access.

### Deploying with CLI

Use the provided CLI tool to deploy the Elasticsearch cluster:

```bash
platon pulumi up --stack-input elasticsearch-cluster.yaml
```

If no Pulumi module is specified, the CLI uses the default module corresponding to the API resource.

**Note**: Ensure that you have the necessary Kubernetes cluster credentials and that the required operators (
Elasticsearch Operator, Istio, Cert-Manager) are installed in your cluster before deploying.
