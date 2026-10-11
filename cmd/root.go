package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	relversion "sigs.k8s.io/release-utils/version"
)

var (
	Open          bool
	Serve         bool
	ServePort     int
	UserAgent     string
	ProxyString   string
	Cookies       []string
	Mirror        bool
	MaxPages      int
	Delay         time.Duration
	InsecureHosts []string
	Pages         []string

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

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			// Create options from flags
			opts := CloneOptions{
				Serve:         Serve,
				Open:          Open,
				ServePort:     ServePort,
				Cookies:       Cookies,
				Proxy:         ProxyString,
				UserAgent:     UserAgent,
				Mirror:        Mirror,
				MaxPages:      MaxPages,
				Delay:         Delay,
				InsecureHosts: InsecureHosts,
				Pages:         Pages,
			}

			// Otherwise.. clone ahead!
			if err := CloneSite(ctx, args, opts); err != nil {
				log.Fatalf("%+v", err)
			}
		},
	}
)

// Execute the clone command
func Execute() {
	// Subcommand `version` from release-utils
	rootCmd.AddCommand(relversion.Version())

	// Persistent Flags
	pf := rootCmd.PersistentFlags()
	pf.BoolVarP(&Open, "open", "o", false, "Automatically open project in default browser")
	pf.BoolVarP(&Serve, "serve", "s", false, "Serve the generated files using Echo.")
	pf.IntVarP(&ServePort, "servePort", "P", 5000, "Serve port number.")
	pf.StringVarP(&ProxyString, "proxy_string", "p", "", "Proxy connection string. Support http and socks5 https://pkg.go.dev/github.com/gocolly/colly#Collector.SetProxy")
	pf.StringVarP(&UserAgent, "user_agent", "u", "", "Custom User Agent")
	pf.BoolVar(&Mirror, "mirror", false, "Capture linked pages, preserve paths, and download HTML/CSS dependencies")
	pf.IntVar(&MaxPages, "max-pages", 100, "Maximum pages to capture in mirror mode")
	pf.DurationVar(&Delay, "delay", 250*time.Millisecond, "Minimum interval between requests in mirror mode")
	pf.StringSliceVar(&InsecureHosts, "insecure-host", nil, "Skip certificate verification ONLY for listed hosts in mirror mode")
	pf.StringSliceVar(&Pages, "page", nil, "Additional same-origin pages to seed in mirror mode (repeatable)")
	rootCmd.Flags().StringSliceVarP(&Cookies, "cookie", "C", nil, "Pre-set these cookies")

	// Execute the command :)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
