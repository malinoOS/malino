package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getKernel() error {
	mainURL := "https://kernel.ubuntu.com/mainline/"

	// Fetch the main page
	res, err := http.Get(mainURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", res.StatusCode)
	}

	// Parse the main page
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	// Find all <a> elements and store them in a slice
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		links = append(links, link)
	})

	// Iterate through the links in reverse order
	for i := len(links) - 1; i >= 0; i-- {
		linkURL := mainURL + links[i]

		// Fetch the link page
		res, err = http.Get(linkURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("status code %d", res.StatusCode)
		}

		// Parse the link page
		doc, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return err
		}

		// Check if the page contains the text "Test amd64 missing"
		pageHTML, _ := doc.Html()
		if !strings.Contains(pageHTML, "Test amd64 missing") {
			// If the page does not contain the text, process the link and stop the loop
			fourthLink, _ := doc.Find("a").Eq(3).Attr("href") // get the fourth link
			fourthLinkURL := linkURL + fourthLink

			// least annoying error handling in go
			if err := downloadFile(fourthLinkURL, "kernel.deb"); err != nil {
				return err
			}

			if err := extractWith7z("kernel.deb"); err != nil {
				return err
			}

			if err := extractWith7z("data.tar"); err != nil {
				return err
			}

			if err := moveBootContentToVmlinuz(); err != nil {
				return err
			}

			if err := os.RemoveAll("boot"); err != nil {
				return err
			}

			if err := os.RemoveAll("usr"); err != nil {
				return err
			}

			if err := os.Remove("control.tar"); err != nil {
				return err
			}

			if err := os.Remove("data.tar"); err != nil {
				return err
			}

			if err := os.Remove("kernel.deb"); err != nil {
				return err
			}

			fmt.Println("Done.")
			return nil
		}
	}

	return fmt.Errorf("no suitable link found")
}

func moveBootContentToVmlinuz() error {
	srcPath := "./boot/"
	destPath := "/home/" + currentUser.Username + "/.malino/vmlinuz"

	srcFiles, err := os.ReadDir(srcPath)
	if err != nil {
		return err
	}

	for _, srcFile := range srcFiles {
		srcFilePath := filepath.Join(srcPath, srcFile.Name())

		srcFileData, err := os.ReadFile(srcFilePath)
		if err != nil {
			return err
		}

		if err := os.WriteFile(destPath, srcFileData, 0644); err != nil {
			return err
		}
	}

	return nil
}
