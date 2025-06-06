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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
	schedulingv1alpha1 "sigs.k8s.io/scheduler-plugins/apis/scheduling/v1alpha1"
)

// ElasticQuotaLister helps list ElasticQuotas.
// All objects returned here must be treated as read-only.
type ElasticQuotaLister interface {
	// List lists all ElasticQuotas in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*schedulingv1alpha1.ElasticQuota, err error)
	// ElasticQuotas returns an object that can list and get ElasticQuotas.
	ElasticQuotas(namespace string) ElasticQuotaNamespaceLister
	ElasticQuotaListerExpansion
}

// elasticQuotaLister implements the ElasticQuotaLister interface.
type elasticQuotaLister struct {
	listers.ResourceIndexer[*schedulingv1alpha1.ElasticQuota]
}

// NewElasticQuotaLister returns a new ElasticQuotaLister.
func NewElasticQuotaLister(indexer cache.Indexer) ElasticQuotaLister {
	return &elasticQuotaLister{listers.New[*schedulingv1alpha1.ElasticQuota](indexer, schedulingv1alpha1.Resource("elasticquota"))}
}

// ElasticQuotas returns an object that can list and get ElasticQuotas.
func (s *elasticQuotaLister) ElasticQuotas(namespace string) ElasticQuotaNamespaceLister {
	return elasticQuotaNamespaceLister{listers.NewNamespaced[*schedulingv1alpha1.ElasticQuota](s.ResourceIndexer, namespace)}
}

// ElasticQuotaNamespaceLister helps list and get ElasticQuotas.
// All objects returned here must be treated as read-only.
type ElasticQuotaNamespaceLister interface {
	// List lists all ElasticQuotas in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*schedulingv1alpha1.ElasticQuota, err error)
	// Get retrieves the ElasticQuota from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*schedulingv1alpha1.ElasticQuota, error)
	ElasticQuotaNamespaceListerExpansion
}

// elasticQuotaNamespaceLister implements the ElasticQuotaNamespaceLister
// interface.
type elasticQuotaNamespaceLister struct {
	listers.ResourceIndexer[*schedulingv1alpha1.ElasticQuota]
}
