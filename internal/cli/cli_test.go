package cli

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRootCmd(t *testing.T) {

	tcs := []struct {
		name     string
		setFlags func(*cobra.Command)
		wantErr  bool
	}{
		{
			name: "happy path",
			setFlags: func(cmd *cobra.Command) {
				cmd.SetArgs([]string{"--file", "stdout"})
			},
			wantErr: false,
		},
		{
			name: "missing required flag",
			setFlags: func(cmd *cobra.Command) {
				cmd.SetArgs([]string{})
			},
			wantErr: true,
		},
		{
			name: "passing invalid delivery times",
			setFlags: func(cmd *cobra.Command) {
				cmd.SetArgs([]string{"--file", "stdout", "--to", "13PM"})
			},
			wantErr: true,
		},
		{
			name: "passing all the flags",
			setFlags: func(cmd *cobra.Command) {
				cmd.SetArgs([]string{"--file", "stdout", "-l", "-p", "10245", "--from", "1PM",
					"--to", "6PM", "-o", "stdout"})
			},
			wantErr: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewRootCmd(func(cmd *cobra.Command, args []string) error {
				return nil
			})
			tc.setFlags(cmd)

			if tc.wantErr {
				assert.NotNil(t, cmd.Execute())
			} else {
				assert.Nil(t, cmd.Execute())
			}

		})
	}
}

func TestMustCli(t *testing.T) {
	tcs := []struct {
		name   string
		cmd    *cobra.Command
		panics bool
	}{
		{
			name: "happy path",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			},
			panics: false,
		},
		{
			name: "happy path",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return errors.New("oh noes")
				},
			},
			panics: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				assert.Panics(t, func() {
					MustCli(tc.cmd)
				})
			} else {
				assert.NotPanics(t, func() {
					MustCli(tc.cmd)
				})
			}
		})
	}
}
