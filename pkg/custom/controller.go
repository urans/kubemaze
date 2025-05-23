package custom

import (
	"context"
	"time"

	apiruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
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

	syncHandler func(ctx context.Context, key string) error
}

// NewLiteFinalizerController creates a new LiteFinalizerController instance
func NewLiteFinalizerController(
	kubeClient clientset.Interface,
	svcInformer coreinformer.ServiceInformer,
) *LiteFinalizerController {

	// 增加 informer 监听事件
	// hpaInformer.Informer().AddEventHandlerWithResyncPeriod(
	// 	cache.ResourceEventHandlerFuncs{
	// 		// 新增 HPA 回调方法
	// 		AddFunc: hpaController.enqueueHPA,
	// 		// 修改 HPA 回调方法
	// 		UpdateFunc: hpaController.updateHPA,
	// 		// 删除 HPA 回调方法
	// 		DeleteFunc: hpaController.deleteHPA,
	// 	},
	// 	// 监听事件的定时周期
	// 	resyncPeriod,
	// )

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
	defer lc.queue.ShutDown()

	for range workers {
		go wait.UntilWithContext(ctx, lc.worker, time.Second)
	}
	<-ctx.Done()
}

func (lc *LiteFinalizerController) worker(ctx context.Context) {
	for lc.processNextWorkItem(ctx) {
	}
}

func (lc *LiteFinalizerController) processNextWorkItem(ctx context.Context) bool {
	key, quit := lc.queue.Get()
	if quit {
		return false
	}
	defer lc.queue.Done(key)

	err := lc.syncHandler(ctx, key.(string))
	lc.handleErr(ctx, err, key)
	return true
}

func (lc *LiteFinalizerController) handleErr(ctx context.Context, err error, key any) {
	// TODO
}
