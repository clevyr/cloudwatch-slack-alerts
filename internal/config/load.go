package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (c *Config) Load(cmd *cobra.Command) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			key := flagToEnv(f.Name)
			if val, ok := os.LookupEnv(key); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
				}
			}
		}
	})
	return errors.Join(errs...)
}

func flagToEnv(f string) string {
	return strings.ToUpper(strings.ReplaceAll(f, "-", "_"))
}
