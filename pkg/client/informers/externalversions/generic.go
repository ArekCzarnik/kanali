// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// This file was automatically generated by informer-gen

package externalversions

import (
	"fmt"

	v2 "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=kanali.io, Version=v2
	case v2.SchemeGroupVersion.WithResource("apikeys"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Kanali().V2().ApiKeys().Informer()}, nil
	case v2.SchemeGroupVersion.WithResource("apikeybindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Kanali().V2().ApiKeyBindings().Informer()}, nil
	case v2.SchemeGroupVersion.WithResource("apiproxies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Kanali().V2().ApiProxies().Informer()}, nil
	case v2.SchemeGroupVersion.WithResource("mocktargets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Kanali().V2().MockTargets().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
