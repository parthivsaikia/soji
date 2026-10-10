package cli

import (
	"fmt"
	"os"

	"github.com/parthivsaikia/soji/internal/cli/plan"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "soji",
	Short: "Clear up container registries by removing unwanted images",
	Long:  "Soji clears up container registries by removing images that are no longer relevant.",
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func init() {
	rootCmd.AddCommand(plan.NewCommand())
}
