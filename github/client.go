package github

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	GIT_INFRA_OWNER    = "jaehong21"
	GIT_INFRA_REPO     = "jaehong21-infra-private"
	GIT_INFRA_PATH     = "route53.tf"
	GIT_CREATE_MESSAGE = "🤖 automate: created route53 resource named "
	GIT_DELETE_MESSAGE = "🤖 automate: deleted route53 resource named "
)

func CreateRoute53(dnsName string) error {
	ctx := context.TODO()
	client := getGithubClient(ctx)

	content, fileContent, err := readTerraformContent(ctx)
	if err != nil {
		return err
	}

	names := parseTerraformContent(content)
	for _, name := range names {
		if name == dnsName {
			return errors.New("resource already exists")
		}
	}

	updatedContent := content + createTerraformResource(dnsName)
	opt := &github.RepositoryContentFileOptions{
		Message: github.String(GIT_CREATE_MESSAGE + dnsName),
		Content: []byte(updatedContent),
		SHA:     fileContent.SHA,
	}

	_, _, err = client.Repositories.UpdateFile(
		ctx, GIT_INFRA_OWNER, GIT_INFRA_REPO, GIT_INFRA_PATH, opt)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoute53(resourceName string) error {
	ctx := context.TODO()
	client := getGithubClient(ctx)

	content, fileContent, err := readTerraformContent(ctx)
	if err != nil {
		return err
	}

	names := parseTerraformContent(content)
	found := false
	for _, name := range names {
		if name == resourceName {
			found = true
			break
		}
	}
	if !found {
		return errors.New("resource not found")
	}

	modifiedContent := deleteTerraformResource(content, resourceName)
	opt := &github.RepositoryContentFileOptions{
		Message: github.String(GIT_DELETE_MESSAGE + resourceName),
		Content: []byte(modifiedContent),
		SHA:     fileContent.SHA,
	}

	_, _, err = client.Repositories.UpdateFile(
		ctx, GIT_INFRA_OWNER, GIT_INFRA_REPO, GIT_INFRA_PATH, opt)
	if err != nil {
		return err
	}

	return nil
}

func getGithubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func readTerraformContent(ctx context.Context) (string, *github.RepositoryContent, error) {
	client := getGithubClient(ctx)

	fileContent, _, _, err := client.Repositories.GetContents(
		ctx, GIT_INFRA_OWNER, GIT_INFRA_REPO, GIT_INFRA_PATH, nil)
	if err != nil {
		return "", nil, err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", nil, err
	}

	return content, fileContent, nil
}
