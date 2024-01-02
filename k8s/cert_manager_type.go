package k8s

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

const CLUSTER_ISSUER_NAME = "letsencrypt-prod"

type Certificate struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type Spec struct {
	SecretName string    `yaml:"secretName"`
	IssuerRef  IssuerRef `yaml:"issuerRef"`
	CommonName string    `yaml:"commonName"`
	DNSNames   []string  `yaml:"dnsNames"`
}

type IssuerRef struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
}

func NewCertificateYAML(namespace, certName string) []byte {
	commonName := "*.jaehong21.com"

	cert := Certificate{
		ApiVersion: "cert-manager.io/v1",
		Kind:       "Certificate",
		Metadata: Metadata{
			Name:      certName,
			Namespace: namespace,
		},
		Spec: Spec{
			SecretName: certName,
			IssuerRef: IssuerRef{
				Name: CLUSTER_ISSUER_NAME,
				Kind: "ClusterIssuer",
			},
			CommonName: commonName,
			DNSNames:   []string{commonName},
		},
	}

	yamlBytes, err := yaml.Marshal(cert)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return yamlBytes
}
