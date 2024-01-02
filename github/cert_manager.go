package github

import "github.com/jaehong21/icarus-proxy/k8s"

func CreateCertificateManifest(namespace, certName string) error {
	certFileContent := k8s.NewCertificateYAML(namespace, certName)

	// NOTE: GitOps folder name must be same as namespace name
	filePath := namespace + "/" + certName + "-cert.yaml"
	if err := createNewFileInRepo(filePath, certFileContent, CERT_CREATE_MESSAGE+certName); err != nil {
		return err
	}

	return nil
}
