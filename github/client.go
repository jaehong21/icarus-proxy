package github

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	GIT_OWNER      = "jaehong21"
	GIT_INFRA_REPO = "jaehong21-infra-private"

	GIT_INFRA_ROUTE53_PATH = "route53.tf"
	ROUTE53_CREATE_MESSAGE = "🤖 automate: created route53 resource named "
	ROUTE53_DELETE_MESSAGE = "🤖 automate: deleted route53 resource named "

	GIT_OPS_REPO = "icarus-gitops"

	CERT_CREATE_MESSAGE = "🤖 automate: created certificate named "
)

func getGithubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func updateFileInRepo(ctx context.Context, client *github.Client, fileContent *github.RepositoryContent, filePath string, updatedContent []byte, commitMessage string) error {
	opt := &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: []byte(updatedContent),
		SHA:     fileContent.SHA,
	}

	// NOTE: Update the file in the Terraform repository
	_, _, err := client.Repositories.UpdateFile(
		ctx, GIT_OWNER, GIT_INFRA_REPO, filePath, opt)
	if err != nil {
		return err
	}

	return nil
}

func createNewFileInRepo(filePath string, fileContent []byte, commitMessage string) error {
	ctx := context.TODO()
	client := getGithubClient(ctx)

	// Prepare the options for creating a new file
	opt := &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: fileContent,
	}

	// NOTE: Create the new file in the GitOps repository
	_, _, err := client.Repositories.CreateFile(
		ctx, GIT_OWNER, GIT_OPS_REPO, filePath, opt)
	if err != nil {
		return err
	}

	return nil
}

func readFileInRepo(ctx context.Context, client *github.Client, path string) (string, *github.RepositoryContent, error) {
	fileContent, _, _, err := client.Repositories.GetContents(
		ctx, GIT_OWNER, GIT_INFRA_REPO, path, nil)
	if err != nil {
		return "", nil, err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", nil, err
	}

	return content, fileContent, nil
}
