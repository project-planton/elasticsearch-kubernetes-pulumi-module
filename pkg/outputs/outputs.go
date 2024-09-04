package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	Namespace                            = "namespace"
	ElasticUsername                      = "username"
	ElasticPasswordSecretName            = "password-secret-name"
	ElasticsearchService                 = "elasticsearch-service"
	ElasticsearchPortForwardCommand      = "elasticsearch-port-forward-command"
	ElasticsearchKubeEndpoint            = "elasticsearch-kube-endpoint"
	ElasticsearchIngressExternalHostname = "elasticsearch-ingress-external-hostname"
	ElasticsearchIngressInternalHostname = "elasticsearch-ingress-internal-hostname"

	KibanaService                 = "kibana-service"
	KibanaPortForwardCommand      = "kibana-port-forward-command"
	KibanaKubeEndpoint            = "kibana-kube-endpoint"
	KibanaIngressExternalHostname = "kibana-ingress-external-hostname"
	KibanaIngressInternalHostname = "kibana-ingress-internal-hostname"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *elasticsearchkubernetes.ElasticsearchKubernetesStackInput) *elasticsearchkubernetes.ElasticsearchKubernetesStackOutputs {
	return &elasticsearchkubernetes.ElasticsearchKubernetesStackOutputs{
		Namespace:                       autoapistackoutput.GetVal(pulumiOutputs, Namespace),
		ElasticUsername:                 autoapistackoutput.GetVal(pulumiOutputs, ElasticUsername),
		ElasticPasswordSecretName:       autoapistackoutput.GetVal(pulumiOutputs, ElasticPasswordSecretName),
		ElasticsearchService:            autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchService),
		ElasticsearchPortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchPortForwardCommand),
		ElasticsearchKubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchKubeEndpoint),
		ElasticsearchExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchIngressExternalHostname),
		ElasticsearchInternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchIngressInternalHostname),
		KibanaService:                   autoapistackoutput.GetVal(pulumiOutputs, KibanaService),
		KibanaPortForwardCommand:        autoapistackoutput.GetVal(pulumiOutputs, KibanaPortForwardCommand),
		KibanaKubeEndpoint:              autoapistackoutput.GetVal(pulumiOutputs, KibanaKubeEndpoint),
		KibanaExternalHostname:          autoapistackoutput.GetVal(pulumiOutputs, KibanaIngressExternalHostname),
		KibanaInternalHostname:          autoapistackoutput.GetVal(pulumiOutputs, KibanaIngressInternalHostname),
	}
}
