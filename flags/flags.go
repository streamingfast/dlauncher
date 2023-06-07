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
	out := make(map[string]bool)
	return recurseCommands(root, nil, out)
}

// ReboundFlagAnnotation has been explicitely choosen to be compatible with `github.com/streamingfast/cli`
var ReboundFlagAnnotation = "github.com/streamingfast/cli#rebound-key"

func recurseCommands(root *cobra.Command, segments []string, flags map[string]bool) map[string]bool {
	var segmentPrefix string
	var segmentPrefixDot string

	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
		segmentPrefixDot = strings.Join(segments, ".") + "."
	}

	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "global-" + f.Name
		newVarDot := segmentPrefixDot + "global." + f.Name

		flags[newVar] = true
		viper.BindPFlag(newVar, f)
		viper.BindPFlag(newVarDot, f)

		addAnnotation(f, ReboundFlagAnnotation, newVarDot)
	})
	root.Flags().VisitAll(func(f *pflag.Flag) {
		newVar := f.Name

		flags[newVar] = true
		viper.BindPFlag(newVar, f)

		addAnnotation(f, ReboundFlagAnnotation, newVar)
	})

	for _, cmd := range root.Commands() {
		flags = recurseCommands(cmd, append(segments, cmd.Name()), flags)
	}
	return flags
}

func addAnnotation(flag *pflag.Flag, key string, value string) {
	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}

	flag.Annotations[key] = append(flag.Annotations[key], value)
}
