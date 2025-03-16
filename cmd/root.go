package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var rootCmd = &cobra.Command{
	Use:   "cli_web_scraper [url]",
	Short: "Scrapes the given domain for links and checks them",
	Long: `Scrapes the given domain for links by crawling through links extracted with html parser.
	While going through them, links get organized into different slices, marking them as either inside, outside or dead.`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run:                   findDead,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

func getPageContent(link string) string {

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("Stopped after 10 redirects")
			}
			return nil
		},
	}

	res, err := client.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(content)

}

func parseAnchors(startingContent string) []string {
	doc, err := html.Parse(strings.NewReader(startingContent))
	if err != nil {
		log.Fatal(err)
	}

	var anchors []string

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					anchors = append(anchors, a.Val)
					break
				}
			}
		}
	}

	return anchors
}

func hasHTTP(anchor string) bool {
	hasHTTP := strings.Index(anchor, "http") == 0
	if hasHTTP {
		return true
	} else {
		return false
	}
}

func makeFullLink(originalLink string, path string) (fullLink string) {
	fullLink = originalLink + path
	return fullLink
}

func isDead(link string) bool {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(link)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}

		log.Fatal(err)
	}
	res.Body.Close()

	if res.StatusCode < 400 {
		return false
	}

	return true
}

func repeatChecks(linksInsideDomain *[]string, linksOutsideDomain *[]string, originalLink string, linksDead *[]string) {

	var newLinkFound bool

	for _, a := range *linksInsideDomain {
		content := getPageContent(a)
		anchors := parseAnchors(content)
		for _, a := range anchors {
			if !hasHTTP(a) {
				n := makeFullLink(originalLink, a)
				if !isDead(n) {
					if !slices.Contains(*linksInsideDomain, n) {
						*linksInsideDomain = append(*linksInsideDomain, n)
						newLinkFound = true
						fmt.Println("Found new inside link - ", n)
					}
				} else {
					if !slices.Contains(*linksDead, n) {
						*linksDead = append(*linksDead, n)
						fmt.Println("Found new dead link - ", n)
					}
				}

			} else {
				if !isDead(a) {
					if !slices.Contains(*linksOutsideDomain, a) {
						*linksOutsideDomain = append(*linksOutsideDomain, a)
						newLinkFound = true
						fmt.Println("Found new outside link - ", a)
					}
				} else {
					if !slices.Contains(*linksDead, a) {
						*linksDead = append(*linksDead, a)
						fmt.Println("Found new dead link - ", a)
					}
				}
			}
		}
	}

	if newLinkFound {
		repeatChecks(linksInsideDomain, linksOutsideDomain, originalLink, linksDead)
	}

}

func findDead(cmd *cobra.Command, args []string) {

	startingLink := args[0]
	startingContent := getPageContent(startingLink)

	fmt.Println("\nStarting link: ", startingLink, "\n")

	var linksInsideDomain []string
	var linksOutsideDomain []string
	var linksDead []string

	for _, a := range parseAnchors(startingContent) {
		if !hasHTTP(a) {
			n := makeFullLink(startingLink, a)
			if !isDead(n) {
				if !slices.Contains(linksInsideDomain, n) {
					linksInsideDomain = append(linksInsideDomain, n)
					fmt.Println("Found new inside link - ", n)
				}
			} else {
				if !slices.Contains(linksDead, n) {
					linksDead = append(linksDead, n)
					fmt.Println("Found new dead link - ", n)
				}
			}
		} else {
			if !isDead(a) {
				if !slices.Contains(linksOutsideDomain, a) {
					linksOutsideDomain = append(linksOutsideDomain, a)
					fmt.Println("Found new outside link - ", a)
				}
			} else {
				if !slices.Contains(linksDead, a) {
					linksDead = append(linksDead, a)
					fmt.Println("Found new dead link - ", a)
				}
			}
		}
	}

	repeatChecks(&linksInsideDomain, &linksOutsideDomain, startingLink, &linksDead)

	fmt.Println("\n\nLinks inside domain: ")
	for i, a := range linksInsideDomain {
		fmt.Println(i+1, " - ", a)
	}

	fmt.Println("\nLinks outside domain: ")
	for i, a := range linksOutsideDomain {
		fmt.Println(i+1, " - ", a)
	}

	fmt.Println("\nDead links: ")
	for i, a := range linksDead {
		fmt.Println(i+1, " - DEAD - ", a)
	}

}
