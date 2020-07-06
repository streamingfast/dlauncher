package flags

import (
	"fmt"
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
	out := make(map[string]bool)
	return recurseCommands(root, nil, out)
}

func recurseCommands(root *cobra.Command, segments []string, flags map[string]bool) map[string]bool {
	var segmentPrefix string
	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
	}

	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "global-" + f.Name
		fmt.Println(newVar)
		flags[newVar] = true
		viper.BindPFlag(newVar, f)
	})
	root.Flags().VisitAll(func(f *pflag.Flag) {
		newVar := f.Name
		fmt.Println(newVar)
		flags[newVar] = true
		viper.BindPFlag(newVar, f)
	})

	for _, cmd := range root.Commands() {
		flags = recurseCommands(cmd, append(segments, cmd.Name()), flags)
	}
	return flags
}
