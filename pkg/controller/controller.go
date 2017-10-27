package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	r2d4com "github.com/r2d4/crd/pkg/apis/r2d4.com"
	crv1 "github.com/r2d4/crd/pkg/apis/r2d4.com/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type R2d4Controller struct {
	R2d4Client rest.Interface
}

// Run starts an Foo resource controller
func (c *R2d4Controller) Run(ctx context.Context) error {
	fmt.Print("Watch Example objects\n")

	// Watch Foo objects
	_, err := c.watchFoos(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Foo resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *R2d4Controller) watchFoos(ctx context.Context) (cache.Controller, error) {
	fmt.Println(c.R2d4Client.APIVersion())
	source := cache.NewListWatchFromClient(
		c.R2d4Client,
		crv1.FooResourcePlural,
		r2d4com.GroupName,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.Foo{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *R2d4Controller) onAdd(obj interface{}) {
	foo := obj.(*crv1.Foo)
	fmt.Printf("[CONTROLLER] OnAdd %s\n", foo.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
	fooCopy.Status = crv1.FooStatus{
		State:   crv1.A,
		Message: "Successfully processed by controller",
	}

	err := c.R2d4Client.Put().
		Name(foo.ObjectMeta.Name).
		Namespace(foo.ObjectMeta.Namespace).
		Resource(crv1.FooResourcePlural).
		Body(fooCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", fooCopy)
	}
}

func (c *R2d4Controller) onUpdate(oldObj, newObj interface{}) {
	oldFoo := oldObj.(*crv1.Foo)
	newFoo := newObj.(*crv1.Foo)
	fmt.Printf("[CONTROLLER] OnUpdate oldObj: %s\n", oldFoo.ObjectMeta.SelfLink)
	fmt.Printf("[CONTROLLER] OnUpdate newObj: %s\n", newFoo.ObjectMeta.SelfLink)
}

func (c *R2d4Controller) onDelete(obj interface{}) {
	foo := obj.(*crv1.Foo)
	fmt.Printf("[CONTROLLER] OnDelete %s\n", foo.ObjectMeta.SelfLink)
}
