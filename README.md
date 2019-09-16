
mkdir stack-linode; cd stack-linode; git init;
kubebuilder init --domain stack.crossplane.io --owner "The Crossplane Authors" --repo github.com/displague/stack-linode
kubebuilder create api --group linode --version v1alpha1 --kind ObjectStorage


add "github.com/displague/stack-linode/controllers" to imports
https://github.com/kubernetes-sigs/kubebuilder/issues/742


