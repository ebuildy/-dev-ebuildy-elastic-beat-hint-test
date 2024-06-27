package cmd

import (
	"fmt"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewAdhocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test elastic beat hint",
		Long: `A ClI to test elastic beat (filebeat...) hint, for docker or kubernetes. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			hintsKey, err := flags.GetString("kind")

			if err != nil {
				ZapLog.Fatal(err)
			}

			aFlags, err := flags.GetStringToString("annotation")

			if err != nil {
				ZapLog.Fatal(err)
			}

			annotations := make(mapstr.M, 0)

			for k, v := range aFlags {
				key := fmt.Sprintf("%s.%s/%s", AnnotationPrefix, hintsKey, k)

				annotations.Put(key, v)

				ZapLog.Infof("add annotation %s: %s", key, v)
			}

			hints := getHints(annotations)
			config := buildConfig(hints, hintsKey)
			d, err := yaml.Marshal(config)

			if err != nil {
				ZapLog.Fatalf("error: %v", err)
			}

			fmt.Println(string(d))
		},
	}

	cmd.Flags().StringP("kind", "k", "logs", "logs, metrics")
	cmd.Flags().StringToStringP("annotation", "a", map[string]string{}, "annotation (eg: enabled=true)")

	return cmd
}
