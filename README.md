kubebuilder init --domain stack.crossplane.io --owner "The Crossplane Authors" --repo github.com/displague/stack-linode
kubebuilder create api --group linode --version v1alpha1 --kind Instance --resource=true --controller=false
kubebuilder create api --group="linode" --version v1alpha1 --kind Provider --resource=true --controller=false
