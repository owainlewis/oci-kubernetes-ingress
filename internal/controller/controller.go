package controller

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type OracleIngressReconciler struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

var _ reconcile.Reconciler = &OracleIngressReconciler{}

func (r *OracleIngressReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("Reconcile...%s", request)
	return reconcile.Result{}, nil
}

func NewOracleIngressReconciler(client client.Client) *OracleIngressReconciler {
	return &OracleIngressReconciler{client: client}
}
