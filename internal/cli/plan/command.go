package plan

import (
	"fmt"

	"github.com/parthivsaikia/soji/internal/service/plan"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Plan status of container images",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			url, err := cmd.Flags().GetString("url")
			if err != nil {
				return fmt.Errorf("error in getting url flag: %w", err)
			}
			if url == "" {
				return fmt.Errorf("url not provided")
			}
			_, err = cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return fmt.Errorf("error in getting kubeconfig path: %w", err)
			}
			_, err = cmd.Flags().GetBool("private")
			if err != nil {
				return fmt.Errorf("error in getting private flag: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
			private, _ := cmd.Flags().GetBool("private")
			url, _ := cmd.Flags().GetString("url")
			err := plan.Plan(kubeconfig, url, private)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String("kubeconfig", "", "path to kubeconfig file")
	cmd.Flags().String("url", "", "url of the container registry")
	cmd.Flags().Bool("private", false, "include private repositories")
	return cmd
}
