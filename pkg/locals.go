package pkg

import (
	elasticsearchkubernetesv1 "buf.build/gen/go/project-planton/apis/protocolbuffers/go/project/planton/provider/kubernetes/elasticsearchkubernetes/v1"
	"fmt"
	"github.com/project-planton/elasticsearch-kubernetes-pulumi-module/pkg/outputs"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	ElasticsearchIngressExternalHostname string
	ElasticsearchIngressInternalHostname string
	ElasticsearchKubePortForwardCommand  string
	ElasticsearchKubeServiceFqdn         string
	ElasticsearchKubeServiceName         string
	Namespace                            string
	ElasticsearchKubernetes              *elasticsearchkubernetesv1.ElasticsearchKubernetes
	KibanaIngressExternalHostname        string
	KibanaIngressInternalHostname        string
	KibanaKubePortForwardCommand         string
	KibanaKubeServiceFqdn                string
	KibanaKubeServiceName                string
	IngressHostnames                     []string
	IngressCertClusterIssuerName         string
	IngressCertSecretName                string
	Labels                               map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *elasticsearchkubernetesv1.ElasticsearchKubernetesStackInput) *Locals {
	locals := &Locals{}
	//assign value for the local variable to make it available across the module.
	locals.ElasticsearchKubernetes = stackInput.Target

	elasticsearchKubernetes := stackInput.Target

	locals.Labels = map[string]string{
		kuberneteslabelkeys.Environment:  stackInput.Target.Metadata.Env.Id,
		kuberneteslabelkeys.Organization: stackInput.Target.Metadata.Org,
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceId:   stackInput.Target.Metadata.Id,
		kuberneteslabelkeys.ResourceKind: "elasticsearch_kubernetes",
	}

	ctx.Export(outputs.ElasticUsername, pulumi.String("elastic"))
	ctx.Export(outputs.ElasticPasswordSecretName, pulumi.Sprintf("%s-es-elastic-user", elasticsearchKubernetes.Metadata.Name))
	ctx.Export(outputs.ElasticPasswordSecretKey, pulumi.String("elastic"))

	//decide on the namespace
	locals.Namespace = elasticsearchKubernetes.Metadata.Id

	locals.ElasticsearchKubeServiceName = fmt.Sprintf("%s-es-http", elasticsearchKubernetes.Metadata.Name)

	//export kubernetes service name
	ctx.Export(outputs.ElasticsearchService, pulumi.String(locals.ElasticsearchKubeServiceName))

	locals.ElasticsearchKubeServiceFqdn = fmt.Sprintf("%s.%s.svc.cluster.local", locals.ElasticsearchKubeServiceName, locals.Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.ElasticsearchKubeEndpoint, pulumi.String(locals.ElasticsearchKubeServiceFqdn))

	locals.ElasticsearchKubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s %d:%d",
		locals.Namespace, locals.ElasticsearchKubeServiceName, vars.ElasticsearchPort, vars.ElasticsearchPort)

	//export kube-port-forward command
	ctx.Export(outputs.ElasticsearchPortForwardCommand, pulumi.String(locals.ElasticsearchKubePortForwardCommand))

	locals.KibanaKubeServiceName = fmt.Sprintf("%s-kb-http", elasticsearchKubernetes.Metadata.Name)

	//export kubernetes service name
	ctx.Export(outputs.KibanaService, pulumi.String(locals.KibanaKubeServiceName))

	locals.KibanaKubeServiceFqdn = fmt.Sprintf("%s.%s.svc.cluster.local", locals.KibanaKubeServiceName, locals.Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.KibanaKubeEndpoint, pulumi.String(locals.KibanaKubeServiceFqdn))

	locals.KibanaKubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s %d:%d",
		locals.Namespace, locals.KibanaKubeServiceName, vars.KibanaPort, vars.KibanaPort)

	//export kube-port-forward command
	ctx.Export(outputs.KibanaPortForwardCommand, pulumi.String(locals.KibanaKubePortForwardCommand))

	if elasticsearchKubernetes.Spec.Ingress == nil ||
		!elasticsearchKubernetes.Spec.Ingress.IsEnabled ||
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return locals
	}

	locals.ElasticsearchIngressExternalHostname = fmt.Sprintf("%s.%s", elasticsearchKubernetes.Metadata.Id,
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName)

	locals.ElasticsearchIngressInternalHostname = fmt.Sprintf("%s-internal.%s", elasticsearchKubernetes.Metadata.Id,
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName)

	locals.KibanaIngressExternalHostname = fmt.Sprintf("%s-kb.%s", elasticsearchKubernetes.Metadata.Id,
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName)

	locals.KibanaIngressInternalHostname = fmt.Sprintf("%s-kb-internal.%s", elasticsearchKubernetes.Metadata.Id,
		elasticsearchKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHostnames = []string{
		locals.ElasticsearchIngressExternalHostname,
		locals.ElasticsearchIngressInternalHostname,
		locals.KibanaIngressExternalHostname,
		locals.KibanaIngressInternalHostname,
	}

	locals.IngressCertClusterIssuerName = elasticsearchKubernetes.Spec.Ingress.EndpointDomainName

	locals.IngressCertSecretName = elasticsearchKubernetes.Metadata.Id

	//export ingress hostnames
	ctx.Export(outputs.ElasticsearchIngressExternalHostname, pulumi.String(locals.ElasticsearchIngressExternalHostname))
	ctx.Export(outputs.ElasticsearchIngressInternalHostname, pulumi.String(locals.ElasticsearchIngressInternalHostname))
	ctx.Export(outputs.KibanaIngressExternalHostname, pulumi.String(locals.KibanaIngressExternalHostname))
	ctx.Export(outputs.KibanaIngressInternalHostname, pulumi.String(locals.KibanaIngressInternalHostname))

	return locals
}
