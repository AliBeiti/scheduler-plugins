/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	gentype "k8s.io/client-go/gentype"
	v1alpha1 "sigs.k8s.io/scheduler-plugins/apis/scheduling/v1alpha1"
	schedulingv1alpha1 "sigs.k8s.io/scheduler-plugins/pkg/generated/applyconfiguration/scheduling/v1alpha1"
	typedschedulingv1alpha1 "sigs.k8s.io/scheduler-plugins/pkg/generated/clientset/versioned/typed/scheduling/v1alpha1"
)

// fakeElasticQuotas implements ElasticQuotaInterface
type fakeElasticQuotas struct {
	*gentype.FakeClientWithListAndApply[*v1alpha1.ElasticQuota, *v1alpha1.ElasticQuotaList, *schedulingv1alpha1.ElasticQuotaApplyConfiguration]
	Fake *FakeSchedulingV1alpha1
}

func newFakeElasticQuotas(fake *FakeSchedulingV1alpha1, namespace string) typedschedulingv1alpha1.ElasticQuotaInterface {
	return &fakeElasticQuotas{
		gentype.NewFakeClientWithListAndApply[*v1alpha1.ElasticQuota, *v1alpha1.ElasticQuotaList, *schedulingv1alpha1.ElasticQuotaApplyConfiguration](
			fake.Fake,
			namespace,
			v1alpha1.SchemeGroupVersion.WithResource("elasticquotas"),
			v1alpha1.SchemeGroupVersion.WithKind("ElasticQuota"),
			func() *v1alpha1.ElasticQuota { return &v1alpha1.ElasticQuota{} },
			func() *v1alpha1.ElasticQuotaList { return &v1alpha1.ElasticQuotaList{} },
			func(dst, src *v1alpha1.ElasticQuotaList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.ElasticQuotaList) []*v1alpha1.ElasticQuota {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.ElasticQuotaList, items []*v1alpha1.ElasticQuota) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
