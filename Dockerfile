FROM alpine:latest

MAINTAINER "Owain Lewis <owain.lewis@oracle.com>"

RUN apk --no-cache add ca-certificates

COPY --from=0 bin/oci-kubernetes-ingress /usr/local/bin/oci-kubernetes-ingress

ENTRYPOINT ["/usr/local/bin/oci-kubernetes-ingress"]
