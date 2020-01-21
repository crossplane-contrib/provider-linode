module github.com/displague/stack-linode

go 1.13

require (
	github.com/crossplaneio/crossplane-runtime v0.4.0
	github.com/crossplaneio/crossplane-tools v0.0.0-20191220202319-9033bd8a02ce
	github.com/linode/linodego v0.12.2
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver v0.17.1
	k8s.io/apimachinery v0.17.1
	k8s.io/client-go v0.17.1
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/controller-tools v0.2.4
)

replace github.com/linodego/linode => /home/marques/.local/share/go/src/github.com/linode/linodego
