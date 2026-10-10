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
		RunE: func(cmd *cobra.Command, args []string) error {
			images, err := plan.Plan()
			if err != nil {
				return fmt.Errorf("error in planning: %w", err)
			}
			for _, image := range images {
				fmt.Println(image)
			}
			return nil
		},
	}
	return cmd
}
