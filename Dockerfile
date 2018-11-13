FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 bin/oci-kubernetes-ingress /usr/local/bin/oci-kubernetes-ingress

ENTRYPOINT ["/usr/local/bin/oci-kubernetes-ingress"]
