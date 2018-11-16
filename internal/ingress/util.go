package ingress

import (
	"fmt"

	"k8s.io/api/extensions/v1beta1"
)

// GetLoadBalancerUniqueName returns a unique name for a given ingress that
// can be used to lookup OCI load balaners.
// TODO can we use OCIDs rather than display names which are not unique?
func GetLoadBalancerUniqueName(ingress *v1beta1.Ingress) string {
	return fmt.Sprintf("%s", ingress.UID)
}
