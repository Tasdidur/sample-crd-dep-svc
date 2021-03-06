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

package v1

import (
	"context"
	"time"

	v1 "github.com/Tasdidur/xcrd/pkg/apis/xapi.com/v1"
	scheme "github.com/Tasdidur/xcrd/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// XcrdsGetter has a method to return a XcrdInterface.
// A group's client should implement this interface.
type XcrdsGetter interface {
	Xcrds(namespace string) XcrdInterface
}

// XcrdInterface has methods to work with Xcrd resources.
type XcrdInterface interface {
	Create(ctx context.Context, xcrd *v1.Xcrd, opts metav1.CreateOptions) (*v1.Xcrd, error)
	Update(ctx context.Context, xcrd *v1.Xcrd, opts metav1.UpdateOptions) (*v1.Xcrd, error)
	UpdateStatus(ctx context.Context, xcrd *v1.Xcrd, opts metav1.UpdateOptions) (*v1.Xcrd, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Xcrd, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.XcrdList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Xcrd, err error)
	XcrdExpansion
}

// xcrds implements XcrdInterface
type xcrds struct {
	client rest.Interface
	ns     string
}

// newXcrds returns a Xcrds
func newXcrds(c *XapiV1Client, namespace string) *xcrds {
	return &xcrds{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the xcrd, and returns the corresponding xcrd object, and an error if there is any.
func (c *xcrds) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Xcrd, err error) {
	result = &v1.Xcrd{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("xcrds").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Xcrds that match those selectors.
func (c *xcrds) List(ctx context.Context, opts metav1.ListOptions) (result *v1.XcrdList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.XcrdList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("xcrds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested xcrds.
func (c *xcrds) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("xcrds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a xcrd and creates it.  Returns the server's representation of the xcrd, and an error, if there is any.
func (c *xcrds) Create(ctx context.Context, xcrd *v1.Xcrd, opts metav1.CreateOptions) (result *v1.Xcrd, err error) {
	result = &v1.Xcrd{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("xcrds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(xcrd).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a xcrd and updates it. Returns the server's representation of the xcrd, and an error, if there is any.
func (c *xcrds) Update(ctx context.Context, xcrd *v1.Xcrd, opts metav1.UpdateOptions) (result *v1.Xcrd, err error) {
	result = &v1.Xcrd{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("xcrds").
		Name(xcrd.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(xcrd).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *xcrds) UpdateStatus(ctx context.Context, xcrd *v1.Xcrd, opts metav1.UpdateOptions) (result *v1.Xcrd, err error) {
	result = &v1.Xcrd{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("xcrds").
		Name(xcrd.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(xcrd).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the xcrd and deletes it. Returns an error if one occurs.
func (c *xcrds) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("xcrds").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *xcrds) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("xcrds").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched xcrd.
func (c *xcrds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Xcrd, err error) {
	result = &v1.Xcrd{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("xcrds").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
