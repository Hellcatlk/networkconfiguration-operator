package controllers

import (
	"context"

	"github.com/Hellcatlk/networkconfiguration-operator/api/v1alpha1"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/machine"
	"github.com/Hellcatlk/networkconfiguration-operator/pkg/utils/finalizer"
	ctrl "sigs.k8s.io/controller-runtime"
)

const switchFinalizerKey string = "foregroundDeletion"

// noneHandler add finalizers to CR
func (r *SwitchReconciler) noneHandler(ctx context.Context, info *machine.ReconcileInfo, instance interface{}) (machine.StateType, ctrl.Result, error) {
	info.Logger.Info("none")

	i := instance.(*v1alpha1.SwitchPort)

	// Add finalizer
	finalizer.Add(&i.Finalizers, switchPortFinalizerKey)

	return v1alpha1.SwitchNone, ctrl.Result{}, nil
}
