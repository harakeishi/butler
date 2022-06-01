package butler

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getTableInfo() []table {
	var tables []table

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	)
	tc := oauth2.NewClient(ctx, ts)

	// エンタープライズには対応してないのでhttps://pkg.go.dev/github.com/google/go-github@v17.0.0+incompatible/github#NewEnterpriseClient　をよしなに使うことを考える
	client := github.NewClient(tc)

	_, contents, _, _ := client.Repositories.GetContents(ctx, "harakeishi", "butler", "sampledb", nil)

	for _, v := range contents {
		path := v.GetPath()
		if strings.Contains(path, ".md") {
			content, _, _, _ := client.Repositories.GetContents(ctx, "harakeishi", "butler", path, nil)
			if strings.Replace(content.GetName(), ".md", "", -1) == "README" {
				continue
			}
			table := table{name: strings.Replace(content.GetName(), ".md", "", -1)}
			text, _ := content.GetContent()
			table.columns = MarkdownParseTocolumn(text)
			tables = append(tables, table)
		}
	}
	return tables
}

func MarkdownParseTocolumn(text string) []column {
	var result []column
	tmp := strings.Split(text, "#")
	rows := strings.Split(tmp[5], "\n")
	for i, v := range rows {
		if i < 4 {
			continue
		}

		colum := strings.Split(v, "|")
		if len(colum) < 8 {
			continue
		}
		result = append(result, column{
			name:    strings.TrimSpace(colum[1]),
			Type:    strings.TrimSpace(colum[2]),
			comment: strings.TrimSpace(colum[8]),
		})
	}
	return result
}
