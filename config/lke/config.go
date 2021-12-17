package lke

import (
	"encoding/base64"

	"github.com/pkg/errors"

	"github.com/crossplane-contrib/terrajet/pkg/config"
)

const (
	errKubeconfigNotString = "kubeconfig attribute is not a string"
)

// Customize configures individual resources by adding custom ResourceConfigurators.
func Customize(p *config.Provider) {
	p.AddResourceConfigurator("linode_lke_cluster", func(r *config.Resource) {

		// we need to override the default group that terrajet generated for
		// this resource, which would be "github"
		r.ShortGroup = ""

		// this property is generated by the Linode provider and must not be set in this provider
		r.ExternalName = config.IdentifierFromProvider
		r.Sensitive.AdditionalConnectionDetailsFn = DecodeKubeconfig
		r.UseAsync = true
	})
}

// DecodeKubeconfig takes "kubeconfig" attribute and decodes it so that it can
// be used without additional processing.
func DecodeKubeconfig(attr map[string]interface{}) (map[string][]byte, error) {
	if attr["kubeconfig"] == nil {
		return nil, nil
	}
	s, ok := attr["kubeconfig"].(string)
	if !ok {
		return nil, errors.New(errKubeconfigNotString)
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode kubeconfig")
	}
	return map[string][]byte{
		"kubeconfig": decoded,
	}, nil
}
