# Annotations

You can add these Kubernetes annotations to specific Ingress objects to customize their behavior.

!!! tip
    Annotation keys and values can only be strings.
    Other types, such as boolean or numeric values must be quoted,
    i.e. `"true"`, `"false"`, `"100"`.

|Name                       | type |
|---------------------------|------|
|[oci.ingress.kubernetes.io/loadbalancer-shape](#loadbalancer) |string|
|[oci.ingress.kubernetes.io/loadbalancer-shape](#loadbalancer) |string|


# LoadBalancer

You can control basic load balancer properties such as the shape and subnets using annotations.
