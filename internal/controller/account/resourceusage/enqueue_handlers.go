package resourceusage

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/sap/crossplane-provider-btp/apis/v1alpha1"
)

type adder interface {
	Add(item any)
}

// EnqueueRequestForResourceUsage enqueues a reconcile.Request for a referenced
// ResourceUsage.
type EnqueueRequestForResourceUsage struct {
}

// Create adds a NamespacedName for the supplied CreateEvent if its Object is a
// ResourceUsageReferencer.
func (e *EnqueueRequestForResourceUsage) Create(ctx context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	addResourceUsage(evt.Object, q)
}

// Update adds a NamespacedName for the supplied UpdateEvent if its Objects are
// a ResourceUsageReferencer.
func (e *EnqueueRequestForResourceUsage) Update(ctx context.Context, evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	addResourceUsage(evt.ObjectOld, q)
	addResourceUsage(evt.ObjectNew, q)
}

// Delete adds a NamespacedName for the supplied DeleteEvent if its Object is a
// ResourceUsageReferencer.
func (e *EnqueueRequestForResourceUsage) Delete(ctx context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	addResourceUsage(evt.Object, q)
}

// Generic adds a NamespacedName for the supplied GenericEvent if its Object is
// a ResourceUsageReferencer.
func (e *EnqueueRequestForResourceUsage) Generic(ctx context.Context, evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	addResourceUsage(evt.Object, q)
}

func addResourceUsage(obj runtime.Object, queue adder) {
	pcr, ok := obj.(*v1alpha1.ResourceUsage)
	if !ok {
		return
	}

	queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{Name: pcr.GetName()}})
}
