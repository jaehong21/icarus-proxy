package github

import (
	"fmt"
	"strings"
)

const (
	TF_ZONE_ID   = "var.jaehong21_com_zone_id"
	TF_DNS_VALUE = "aws_eip.k3s_static_ip[0].public_ip"
)

func parseTerraformContent(content string) []string {
	var names []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				name = strings.Trim(name, "\"")
				names = append(names, name)
			}
		}
	}
	return names
}

// Generates a Terraform resource string for a given name
func createTerraformResource(name string) string {
	return ("\n\n" + fmt.Sprintf(`resource "cloudflare_record" "%s" {
  zone_id = %s
  name    = "%s"
  value   = %s
  proxied = false
  type    = "A"
  ttl     = 1
}`, name, TF_ZONE_ID, name, TF_DNS_VALUE))
}

// Deletes a Terraform resource with the specified name from the content
func deleteTerraformResource(content string, resourceName string) string {
	var result []string
	insideBlock := false

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "resource") && strings.Contains(trimmedLine, resourceName) {
			insideBlock = true
			continue
		}
		if insideBlock && trimmedLine == "}" {
			insideBlock = false
			continue
		}
		if !insideBlock {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
