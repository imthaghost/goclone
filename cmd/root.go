package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Flags
	Login bool
	Serve bool

	// Root cmd
	rootCmd = &cobra.Command{
		Use:   "goclones",
		Short: "Clone a website with ease!",
		Long:  `Copy websites to your computer! goclone is a utility that allows you to download a website from the Internet to a local directory. Get html, css, js, images, and other files from the server to your computer. goclone arranges the original site's relative link-structure. Simply open a page of the "mirrored" website in your browser, and you can browse the site from link to link, as if you were viewing it online.`, // TODO Update link once we change repo name
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Print the usage if no args are passed in :)
			if len(args) < 1 {
				if err := cmd.Usage(); err != nil {
					log.Fatal(err)
				}

				return
			}

			// Otherwise.. clone ahead!
			cloneSite(args)
		},
	}
)

// Execute the clone command
func Execute() {
	// Persistent Flags
	rootCmd.PersistentFlags().BoolVarP(&Login, "login", "l", false, "Wether to use a username or password")
	rootCmd.PersistentFlags().BoolVarP(&Serve, "serve", "s", false, "Serve the generated files using Echo.")

	// Execute the command :)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
