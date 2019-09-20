module github.com/displague/stack-linode

go 1.12

require (
	github.com/crossplaneio/crossplane-runtime v0.0.0-20190919002909-d8050430d1b6
	github.com/go-logr/logr v0.1.0
	github.com/linode/linodego v0.10.0
	github.com/onsi/ginkgo v1.6.0
	github.com/onsi/gomega v1.4.2
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/tools v0.0.0-20190920130846-1081e67f6b77 // indirect
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/gengo v0.0.0-20190907103519-ebc107f98eab // indirect
	sigs.k8s.io/controller-runtime v0.2.0-beta.1
)
