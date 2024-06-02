package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

	// Find the last <a> element
	var lastLink string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		lastLink, _ = s.Attr("href")
	})

	// Make sure the last link is a full URL
	lastLinkURL := mainURL + lastLink

	// Fetch the last link page
	res, err = http.Get(lastLinkURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", res.StatusCode)
	}

	// Parse the last link page
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	// Find the fourth <a> element
	var fourthLink string
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 3 { // since indices are 0-based, 3 is the fourth element
			fourthLink, _ = s.Attr("href")
			return false // stop iteration
		}
		return true // continue iteration
	})

	// Make sure the fourth link is a full URL
	fourthLinkURL := lastLinkURL + fourthLink

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

func moveBootContentToVmlinuz() error {
	srcPath := "./boot/"
	destPath := "./vmlinuz"

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
