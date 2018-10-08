package app

import "github.com/stakater/Xposer/internal/pkg/cmd"

// Run runs the command
func Run() error {
	cmd := cmd.NewXposerCommand()
	return cmd.Execute()
}
