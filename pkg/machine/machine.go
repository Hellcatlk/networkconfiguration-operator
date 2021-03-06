package machine

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StateType is the type of .status.state
type StateType string

// Machine is a state machine
type Machine struct {
	info     *ReconcileInfo
	instance Instance
	handlers map[StateType]Handler
}

// ReconcileInfo is the information need by reconcile
type ReconcileInfo struct {
	Client client.Client
	Logger logr.Logger
}

// Instance is a object for the CR need be reconcile
// NOTE: Instance must be a pointer
type Instance interface {
	runtime.Object
	GetState() StateType
	SetState(state StateType)
	SetError(err error)
}

// Handler is a state handle function. If an error isn't nil or
// ctrl.Result.Requeue is true the state machine will requeue the Request again
type Handler func(ctx context.Context, info *ReconcileInfo, instance interface{}) (StateType, ctrl.Result, error)

// New a state machine
// NOTE: The paramater of instance must be a pointer
func New(info *ReconcileInfo, instance Instance, handlers map[StateType]Handler) Machine {
	return Machine{
		info:     info,
		instance: instance,
		handlers: handlers,
	}
}

// Reconcile state machine. If dirty is true, it means the instance has changed
func (m *Machine) Reconcile(ctx context.Context) (bool, ctrl.Result, error) {
	m.info.Logger.Info(string(m.instance.GetState()))

	// There are any handler in handlers?
	if m.handlers == nil {
		return false, ctrl.Result{}, fmt.Errorf("haven't any handler")
	}

	// Check the state's handler exist or not
	handler, exist := m.handlers[m.instance.GetState()]
	if !exist {
		return false, ctrl.Result{}, fmt.Errorf("no handler for the state(%s)", m.instance.GetState())
	}

	// Call handler
	instanceDeepCopy := m.instance.DeepCopyObject()
	nextState, result, err := handler(ctx, m.info, m.instance)
	if err != nil {
		err = fmt.Errorf("%s state handler error: %s", m.instance.GetState(), err)
	}
	m.instance.SetState(nextState)
	m.instance.SetError(err)

	// Check instance is dirty or not
	return !reflect.DeepEqual(m.instance, instanceDeepCopy), result, err
}
