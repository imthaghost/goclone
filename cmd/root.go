package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	Open        bool
	Serve       bool
	ServePort   int
	UserAgent   string
	ProxyString string
	Referer     string
	Depth       int
	Cookies     []string

	// Root cmd
	rootCmd = &cobra.Command{
		Use:   "goclone <url>",
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
			if Depth == 0 {
				Depth = 1
			}
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()
			// Otherwise.. clone ahead!
			if err := cloneSite(ctx, args); err != nil {
				log.Fatalf("%+v", err)
			}
		},
	}
)

// Execute the clone command
func Execute() {
	// Persistent Flags
	pf := rootCmd.PersistentFlags()
	pf.BoolVarP(&Open, "open", "o", false, "Automatically open project in deafult browser")
	// rootCmd.PersistentFlags().BoolVarP(&Login, "login", "l", false, "Wether to use a username or password")
	pf.BoolVarP(&Serve, "serve", "s", false, "Serve the generated files using Echo.")
	pf.IntVarP(&ServePort, "servePort", "P", 5000, "Serve port number.")
	pf.IntVarP(&Depth, "depth", "d", 5000, "max depth")
	pf.StringVarP(&ProxyString, "proxy_string", "p", "", "Proxy connection string. Support http and socks5 https://pkg.go.dev/github.com/gocolly/colly#Collector.SetProxy")
	pf.StringVarP(&UserAgent, "user_agent", "u", "", "Custom User Agent")
	pf.StringVarP(&Referer, "referer", "r", "", "Custom Referer")
	rootCmd.Flags().StringSliceVarP(&Cookies, "cookie", "C", nil, "Pre-set these cookies")

	// Execute the command :)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
