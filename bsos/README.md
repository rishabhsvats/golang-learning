
Building bsos 
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"
```

Related Plugins:
- https://github.com/kubernetes-csi/external-provisioner
- https://github.com/kubernetes-csi/external-attacher
- https://github.com/kubernetes-csi/node-driver-registrar
- https://github.com/kubernetes/design-proposals-archive/blob/main/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes
- https://kubernetes.io/blog/2019/01/15/container-storage-interface-ga/#how-to-write-a-csi-driver
- https://github.com/container-storage-interface/spec/blob/master/spec.md#architecture