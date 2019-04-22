package clusterenv

import (
	"context"

	envv1alpha1 "github.com/jmckind/env-operator/pkg/apis/env/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_clusterenv")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ClusterEnv Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileClusterEnv{client: mgr.GetClient(), config: mgr.GetConfig(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("clusterenv-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ClusterEnv
	err = c.Watch(&source.Kind{Type: &envv1alpha1.ClusterEnv{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ClusterEnv
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &envv1alpha1.ClusterEnv{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileClusterEnv{}

// ReconcileClusterEnv reconciles a ClusterEnv object
type ReconcileClusterEnv struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	config *rest.Config
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ClusterEnv object and makes changes based on the state read
// and what is in the ClusterEnv.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileClusterEnv) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ClusterEnv")

	// Fetch the ClusterEnv instance
	instance := &envv1alpha1.ClusterEnv{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	product := "Kubernetes"

	kubeClient, err := kubernetes.NewForConfig(r.config)
	if err != nil {
		reqLogger.Error(err, "unable to create kubeclient")
		return reconcile.Result{}, err
	}

	err = discovery.ServerSupportsVersion(kubeClient.Discovery(), schema.GroupVersion{
		Group:   "config.openshift.io",
		Version: "v1",
	})

	if err == nil {
		// TODO: Fetch ClusterVersion for the actual OpenShift version!
		product = "OpenShift 4.x"
	} else {
		err = discovery.ServerSupportsVersion(kubeClient.Discovery(), schema.GroupVersion{
			Group:   "route.openshift.io",
			Version: "v1",
		})

		if err == nil {
			// TODO: Determine exact OpenShift 3 version
			product = "OpenShift 3.x"
		}
	}

	if product != instance.Status.Product {
		reqLogger.Info("Updating product...", "product", product)
		instance.Status.Product = product
		r.client.Status().Update(context.TODO(), instance)
	}

	return reconcile.Result{}, nil
}
