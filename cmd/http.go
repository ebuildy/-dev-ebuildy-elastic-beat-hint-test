package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func formatToYAML(a any) string {
	b := bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(&a)

	return strings.TrimSpace(b.String())
}

func handleHTTPQuery(annotationsQuery []string, hintsKey string) map[string]interface{} {
	annotations := make(mapstr.M, 0)

	for _, query := range annotationsQuery {
		components := strings.SplitN(query, "=", 2)
		key := fmt.Sprintf("%s.%s/%s", AnnotationPrefix, hintsKey, components[0])

		annotations.Put(key, components[1])

		ZapLog.Infof("add annotation %s: %s", key, components[1])
	}

	hints := getHints(annotations)
	config := buildConfig(hints, hintsKey)

	return map[string]interface{}{
		"annotations": annotationsQuery,
		"hints":       hints,
		"results":     config,
	}
}

func NewHTTPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Run HTTP API server",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.New()
			r.SetFuncMap(template.FuncMap{
				"toYAML": formatToYAML,
			})

			r.LoadHTMLFiles("./template/index.html")

			r.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})

			r.GET("/", func(c *gin.Context) {
				rawAnnotations := c.DefaultQuery("annotations", "enabled=true\nprocessors.drop_event={\"when\":{\"or\":[{\"equals\": {\"log.level\": \"info\"}}]}}")
				hintsKey := c.DefaultQuery("key", "logs")

				r := handleHTTPQuery(strings.Split(rawAnnotations, "\n"), hintsKey)

				c.HTML(http.StatusOK, "index.html", gin.H{
					"annotations":  rawAnnotations,
					"resultConfig": r["results"],
					"resultHints":  r["hints"],
				})
			})

			r.GET("/test", func(c *gin.Context) {
				hintsKey := c.DefaultQuery("key", "logs")
				annotationsQuery := c.QueryArray("a")

				c.JSON(http.StatusOK, handleHTTPQuery(annotationsQuery, hintsKey))
			})

			r.Run()
		},
	}

	return cmd
}
