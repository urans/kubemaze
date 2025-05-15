package custom

import (
	"context"
	"time"

	apiruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformer "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const liteFinalizer = "kallen.io/lite-finalizer"

// LiteFinalizerController is a controller that implements custom finalizers
type LiteFinalizerController struct {
	kubeClient clientset.Interface

	svcLister corelister.ServiceLister

	svcListerSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface
}

// NewLiteFinalizerController creates a new LiteFinalizerController instance
func NewLiteFinalizerController(
	kubeClient clientset.Interface,
	svcInformer coreinformer.ServiceInformer,
) *LiteFinalizerController {

	svcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			// TODO
		},
		DeleteFunc: func(obj any) {
			// TODO
		},
		UpdateFunc: func(old, new any) {
			// TODO
		},
	})

	lc := &LiteFinalizerController{
		kubeClient:      kubeClient,
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "LiteFinalizerController"),
		svcLister:       svcInformer.Lister(),
		svcListerSynced: svcInformer.Informer().HasSynced,
	}
	return lc
}

// Run starts watching and syncing service finalizers
func (lc *LiteFinalizerController) Run(ctx context.Context, workers int) {
	defer apiruntime.HandleCrash()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):

		}
	}
}
