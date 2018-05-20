# Ingress OCI

An [Ingress][0] Controller for Oracle Cloud Infrastructure. There are two load balancing strategies

1. Nodeport (load balance across all nodes/nodeport pairs)
2. PodIP (load balance across all the pod IPs that form a service object)

See [here][1] for more information about OCI path and virtual host based routing.

[0]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[1]: https://docs.us-phoenix-1.oraclecloud.com/Content/Balance/Tasks/managingrequest.htm
