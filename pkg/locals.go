package pkg

import (
	"fmt"
	"github.com/plantoncloud/elasticsearch-kubernetes-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/apiresource/enums/apiresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	IngressExternalHostname      string
	IngressInternalHostname      string
	KubePortForwardCommand       string
	KubeServiceFqdn              string
	KubeServiceName              string
	Namespace                    string
	ElasticsearchKubernetes      *elasticsearchkubernetes.ElasticsearchKubernetes
	ElasticsearchPodSectorLabels map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *elasticsearchkubernetes.ElasticsearchKubernetesStackInput) *Locals {
	locals := &Locals{}
	//assign value for the local variable to make it available across the module.
	locals.ElasticsearchKubernetes = stackInput.ApiResource

	elasticsearchKubernetes := stackInput.ApiResource

	//decide on the namespace
	locals.Namespace = elasticsearchKubernetes.Metadata.Id

	locals.ElasticsearchPodSectorLabels = map[string]string{
		"planton.cloud/resource-kind": apiresourcekind.ApiResourceKind_elasticsearch_kubernetes.String(),
		"planton.cloud/resource-id":   elasticsearchKubernetes.Metadata.Id,
	}

	locals.KubeServiceName = fmt.Sprintf("%s-master", elasticsearchKubernetes.Metadata.Name)

	//export kubernetes service name
	ctx.Export(outputs.ServiceOutputName, pulumi.String(locals.KubeServiceName))

	locals.KubeServiceFqdn = fmt.Sprintf("%s.%s.svc.cluster.local", locals.KubeServiceName, locals.Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.KubeEndpointOutputName, pulumi.String(locals.KubeServiceFqdn))

	locals.KubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s 8080:8080",
		locals.Namespace, locals.KubeServiceName)

	//export kube-port-forward command
	ctx.Export(outputs.PortForwardCommandOutputName, pulumi.String(locals.KubePortForwardCommand))

	if elasticsearchKubernetes.Spec.Ingress == nil ||
		!elasticsearchKubernetes.Spec.Ingress.IsEnabled ||
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return locals
	}

	locals.IngressExternalHostname = fmt.Sprintf("%s.%s", elasticsearchKubernetes.Metadata.Id,
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressInternalHostname = fmt.Sprintf("%s-internal.%s", elasticsearchKubernetes.Metadata.Id,
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName)

	//export ingress hostnames
	ctx.Export(outputs.IngressExternalHostnameOutputName, pulumi.String(locals.IngressExternalHostname))
	ctx.Export(outputs.IngressInternalHostnameOutputName, pulumi.String(locals.IngressInternalHostname))

	return locals
}
