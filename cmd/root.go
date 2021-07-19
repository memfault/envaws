package cmd

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/memfault/envaws/param_providers"
	"github.com/memfault/envaws/poller"
	"github.com/memfault/envaws/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "envaws",
		Short: "A process monitor that restarts a subprocess upon AWS SSM/S3 config changes",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_cmd *cobra.Command, args []string) {
			var provider param_providers.ParamProvider
			wantedParams := viper.GetStringSlice("params")
			pollingIntervalSeconds := viper.GetInt("interval")
			log.Print("Looking for parameters [", strings.Join(wantedParams, ","), "]")
			switch viper.GetString("service") {
			case "ssm":
				provider = param_providers.NewSSMProvider(wantedParams)
			case "s3":
				provider = param_providers.NewS3Provider()
			case "file":
				provider = param_providers.NewFileProvider("./test.json")
			}

			change_detected_channel := make(chan bool)

			// start watching params from the proper provider
			go poller.Poll(provider, int64(pollingIntervalSeconds), change_detected_channel)

			// start the process we are managing
			cmdString := strings.Join(args, " ")
			command := exec.Command(cmdString)
			go runner.RunCmd(command)

			<-change_detected_channel
			log.Println("Change in environment's config detected, exiting gracefully")
			runner.SoftThenHardKill(command, time.Duration(10)*time.Second)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is to check $HOME/.envaws and $CWD/.envaws)")
	rootCmd.PersistentFlags().StringP("service", "s", "", "whether to use AWS SSM or S3")
	rootCmd.PersistentFlags().IntP("interval", "i", 10, "how often to poll for config changes [in seconds]")

	viper.BindPFlag("service", rootCmd.PersistentFlags().Lookup("service"))
	viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval"))

	viper.SetDefault("service", "ssm") // TODO: validate
	viper.SetDefault("interval", "15") // TODO: validate
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".envaws")
		viper.SetConfigType("yaml")

		// check ~/.envaws and $CWD/.envaws
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
