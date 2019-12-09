module github.com/displague/stack-linode

go 1.12

require (
	cloud.google.com/go v0.38.0 // indirect
	github.com/crossplaneio/crossplane-runtime v0.0.0-20190919002909-d8050430d1b6
	github.com/crossplaneio/crossplane-tools v0.0.0-20191023215726-61fa1eff2a2e
	github.com/linode/linodego v0.10.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/appengine v1.5.0 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.2
	sigs.k8s.io/controller-tools v0.2.4
)

replace github.com/linodego/linode => /home/marques/.local/share/go/src/github.com/linode/linodego
