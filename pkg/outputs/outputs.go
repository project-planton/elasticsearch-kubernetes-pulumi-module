package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/kubernetes"
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
		Namespace: autoapistackoutput.GetVal(pulumiOutputs, Namespace),
		Elasticsearch: &elasticsearchkubernetes.ElasticsearchKubernetesElasticsearchStackOutputs{
			Username: autoapistackoutput.GetVal(pulumiOutputs, ElasticUsername),
			//ElasticPasswordSecretName:       ,
			Service:            autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchService),
			PortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchPortForwardCommand),
			KubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchKubeEndpoint),
			ExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchIngressExternalHostname),
			InternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, ElasticsearchIngressInternalHostname),
			PasswordSecret: &kubernetes.KubernernetesSecretKey{
				Name: autoapistackoutput.GetVal(pulumiOutputs, ElasticPasswordSecretName),
				Key:  "coming-soon",
			},
		},
		Kibana: &elasticsearchkubernetes.ElasticsearchKubernetesKibanaStackOutputs{
			Service:            autoapistackoutput.GetVal(pulumiOutputs, KibanaService),
			PortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, KibanaPortForwardCommand),
			KubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, KibanaKubeEndpoint),
			ExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, KibanaIngressExternalHostname),
			InternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, KibanaIngressInternalHostname),
		},
	}
}
