package pkg

import (
	"github.com/pkg/errors"
	elasticsearchv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/elasticsearch/elasticsearch/v1beta1"
	kibanav1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/elasticsearch/kibana/v1beta1"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/elasticsearchkubernetes"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *elasticsearchkubernetes.ElasticsearchKubernetesStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := initializeLocals(ctx, s.Input)
	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesClusterCredential(ctx,
		s.Input.KubernetesClusterCredential, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	createdNamespace, err := kubernetescorev1.NewNamespace(ctx, locals.ElasticsearchKubernetes.Metadata.Id,
		&kubernetescorev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(locals.ElasticsearchKubernetes.Metadata.Id),
					Labels: pulumi.ToStringMap(s.Labels),
				}),
		}, pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create namespace")
	}

	nodeSets := elasticsearchv1.ElasticsearchSpecNodeSetsArray{}
	nodeSets = append(nodeSets, elasticsearchv1.ElasticsearchSpecNodeSetsArgs{
		Config: pulumi.Map{
			"node.roles":            pulumi.ToStringArray([]string{"master"}),
			"node.store.allow_mmap": pulumi.Bool(false),
		},
		Count: pulumi.Int(1),
		Name:  pulumi.String("master-nodes"),
		PodTemplate: pulumi.Map{
			"metadata": pulumi.Map{
				"labels": pulumi.StringMap{
					"role": pulumi.String("master"),
				},
			},
			"spec": pulumi.Map{
				"containers": pulumi.Array{
					pulumi.Map{
						"name": pulumi.String("elasticsearch"),
						"resources": pulumi.Map{
							"requests": pulumi.Map{
								"memory": pulumi.String("2Gi"),
								"cpu":    pulumi.String("1"),
							},
							"limits": pulumi.Map{
								"memory": pulumi.String("2Gi"),
								"cpu":    pulumi.String("1"),
							},
						},
					},
				},
			},
		},
	})

	nodeSets = append(nodeSets, elasticsearchv1.ElasticsearchSpecNodeSetsArgs{
		Name:  pulumi.String("data-nodes"),
		Count: pulumi.Int(1),
		Config: pulumi.Map{
			"node.roles":            pulumi.ToStringArray([]string{"data"}),
			"node.store.allow_mmap": pulumi.Bool(false),
		},
		PodTemplate: pulumi.Map{
			"metadata": pulumi.Map{
				"labels": pulumi.StringMap{
					"role": pulumi.String("master"),
				},
			},
			"spec": pulumi.Map{
				"containers": pulumi.Array{
					pulumi.Map{
						"name": pulumi.String("elasticsearch"),
						"resources": pulumi.Map{
							"requests": pulumi.Map{
								"memory": pulumi.String("4Gi"),
								"cpu":    pulumi.String("2"),
							},
							"limits": pulumi.Map{
								"memory": pulumi.String("4Gi"),
								"cpu":    pulumi.String("2"),
							},
						},
					},
				},
			},
		},
		VolumeClaimTemplates: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesArray{
			elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesArgs{
				Metadata: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesMetadataArgs{
					Name: pulumi.String("elasticsearch-data"),
				},
				Spec: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesSpecArgs{
					AccessModes: pulumi.StringArray{
						pulumi.String("ReadWriteOnce"),
					},
					Resources: &elasticsearchv1.ElasticsearchSpecNodeSetsVolumeClaimTemplatesSpecResourcesArgs{
						Requests: pulumi.Map{
							"storage": pulumi.String("2Gi"),
						},
					},
				},
			},
		},
	})

	createdElasticSearch, err := elasticsearchv1.NewElasticsearch(ctx, locals.ElasticsearchKubernetes.Metadata.Name, &elasticsearchv1.ElasticsearchArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.ElasticsearchKubernetes.Metadata.Name),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(s.Labels),
			Annotations: pulumi.StringMap{
				"pulumi.com/patchForce": pulumi.String("true"),
			},
		},
		Spec: &elasticsearchv1.ElasticsearchSpecArgs{
			NodeSets: nodeSets,
			Version:  pulumi.String("8.15.1"),
		},
	}, pulumi.Parent(createdNamespace))
	if err != nil {
		return errors.Wrapf(err, "failed to create elastic search")
	}

	_, err = kibanav1.NewKibana(ctx, locals.ElasticsearchKubernetes.Metadata.Id, &kibanav1.KibanaArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.ElasticsearchKubernetes.Metadata.Name),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(s.Labels),
		},
		Spec: &kibanav1.KibanaSpecArgs{
			Version: pulumi.String("8.15.1"),
			Count:   pulumi.Int(1),
			PodTemplate: pulumi.Map{
				"spec": pulumi.Map{
					"containers": pulumi.Array{
						pulumi.Map{
							"name": pulumi.String("kibana"),
							"resources": pulumi.Map{
								"requests": pulumi.Map{
									"memory": pulumi.String("512Mi"),
									"cpu":    pulumi.String("500m"),
								},
								"limits": pulumi.Map{
									"memory": pulumi.String("512Mi"),
									"cpu":    pulumi.String("500m"),
								},
							},
						},
					},
				},
			},
			ElasticsearchRef: kibanav1.KibanaSpecElasticsearchRefArgs{
				Name:      createdElasticSearch.Metadata.Name().Elem(),
				Namespace: createdNamespace.Metadata.Name(),
			},
		},
	}, pulumi.Parent(createdNamespace), pulumi.DependsOn([]pulumi.Resource{createdElasticSearch}))
	if err != nil {
		return errors.Wrapf(err, "failed to create kibana instance for the elastic search instance")
	}

	return nil
}
