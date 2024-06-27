/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/elastic/elastic-agent-autodiscover/utils"
	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AnnotationPrefix = "co.elastic"
var AllSupportedHints = []string{"enabled", "module", "metricsets", "hosts", "period", "timeout", "metrics_path", "username", "password", "stream", "processors", "multiline", "json", "disable", "ssl", "metrics_filters", "raw", "include_lines", "exclude_lines", "fileset", "pipeline", "raw"}

var ZapLog *zap.SugaredLogger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elastic-beat-hint-test",
	Short: "Test elastic beat hint",
	Long: `A ClI to test elastic beat (filebeat...) hint, for docker or kubernetes. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func getHints(annotations mapstr.M) mapstr.M {
	hints, incorrecthints := utils.GenerateHints(annotations, "", AnnotationPrefix, true, AllSupportedHints)

	for _, value := range incorrecthints {
		ZapLog.Infof("provided hint: %s/%s is not in the supported list", AnnotationPrefix, value)
	}

	for _, value := range hints {
		ZapLog.Infof("provided hint: %s/%s is OK list", AnnotationPrefix, value)
	}

	return hints
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
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	enccoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = enccoderConfig

	logger, err := config.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}

	ZapLog = logger.Sugar()

	rootCmd.AddCommand(NewAdhocCommand())
	rootCmd.AddCommand(NewHTTPCommand())
}
