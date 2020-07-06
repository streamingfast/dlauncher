package flags

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func AutoBind(root *cobra.Command, prefix string) map[string]bool {
	viper.SetEnvPrefix(strings.ToUpper(prefix))
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	return recurseCommands(root, nil)
}

func recurseCommands(root *cobra.Command, segments []string) map[string]bool {
	out := make(map[string]bool)
	var segmentPrefix string
	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
	}

	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "global-" + f.Name
		out[newVar] = true
		viper.BindPFlag(newVar, f)
	})
	root.Flags().VisitAll(func(f *pflag.Flag) {
		newVar := f.Name
		out[newVar] = true
		viper.BindPFlag(newVar, f)
	})

	for _, cmd := range root.Commands() {
		recurseCommands(cmd, append(segments, cmd.Name()))
	}
	return out
}
