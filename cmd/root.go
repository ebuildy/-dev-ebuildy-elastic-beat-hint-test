/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/elastic/elastic-agent-autodiscover/utils"
	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var AnnotationPrefix = "co.elastic"
var AllSupportedHints = []string{"enabled", "module", "metricsets", "hosts", "period", "timeout", "metrics_path", "username", "password", "stream", "processors", "multiline", "json", "disable", "ssl", "metrics_filters", "raw", "include_lines", "exclude_lines", "fileset", "pipeline", "raw"}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elastic-beat-hint-test",
	Short: "Test elastic beat hint",
	Long: `A ClI to test elastic beat (filebeat...) hint, for docker or kubernetes. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
		sugar := logger.Sugar()

		flags := cmd.Flags()
		hintsKey, err := flags.GetString("kind")

		if err != nil {
			sugar.Fatal(err)
		}

		aFlags, err := flags.GetStringToString("annotation")

		if err != nil {
			sugar.Fatal(err)
		}

		annotations := make(mapstr.M, 0)

		for k, v := range aFlags {
			key := fmt.Sprintf("%s.%s/%s", AnnotationPrefix, hintsKey, k)

			annotations.Put(key, v)

			sugar.Infof("add annotation %s: %s", key, v)
		}

		hints, incorrecthints := utils.GenerateHints(annotations, "", AnnotationPrefix, true, AllSupportedHints)

		for _, value := range incorrecthints {
			sugar.Infof("provided hint: %s/%s is not in the supported list", AnnotationPrefix, value)
		}

		for _, value := range hints {
			sugar.Infof("provided hint: %s/%s is OK list", AnnotationPrefix, value)
		}

		//annotations = getHintMapStr(annotations, AnnotationPrefix)

		config := buildConfig(hints, hintsKey)

		d, err := yaml.Marshal(config)

		if err != nil {
			sugar.Fatalf("error: %v", err)
		}

		fmt.Println(string(d))
	},
}

/**
* Copy from https://github.com/elastic/beats/blob/main/filebeat/autodiscover/builder/hints/logs.go#L84
 */
func buildConfig(hints mapstr.M, hintsKey string) map[string]any {
	config := make(map[string]any)
	configProcessors := make([]interface{}, 0)
	configIncludeLines := make([]string, 0)
	configExcludeLines := make([]string, 0)
	configJSON := make(map[string]interface{}, 0)

	if procs := utils.GetProcessors(hints, hintsKey); len(procs) != 0 {
		for _, proc := range procs {
			configProcessors = append(configProcessors, proc)
		}
	}

	if json := utils.GetHintMapStr(hints, hintsKey, "json"); len(json) != 0 {
		configJSON = json
	}

	if lines := utils.GetHintAsList(hints, hintsKey, "include_lines"); len(lines) != 0 {
		configIncludeLines = append(configIncludeLines, lines...)
	}

	if lines := utils.GetHintAsList(hints, hintsKey, "exclude_lines"); len(lines) != 0 {
		configExcludeLines = append(configExcludeLines, lines...)
	}

	config["processors"] = configProcessors
	config["includeLines"] = configIncludeLines
	config["excludeLines"] = configExcludeLines
	config["enabled"] = utils.IsEnabled(hints, hintsKey)
	config["json"] = configJSON

	return config
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("kind", "k", "logs", "logs, metrics")
	rootCmd.PersistentFlags().StringToStringP("annotation", "a", map[string]string{}, "annotation (eg: enabled=true)")
}
