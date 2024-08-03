package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getKernel(downloadModules bool) error {
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
			fmt.Println("Downloading kernel...")
			if _, err := os.Stat(fmt.Sprintf("/home/%s/.malino/vmlinuz", currentUser.Username)); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				if err := os.Remove(fmt.Sprintf("/home/%s/.malino/vmlinuz", currentUser.Username)); err != nil {
					return err
				}
			}

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

			if downloadModules {
				fmt.Printf("Linux kernel image \"%s\" downloaded, now downloading kernel modules... (this will take a while)\n", doc.Find("a").Eq(3).Text())

				if _, err := os.Stat(fmt.Sprintf("/home/%s/.malino/modules", currentUser.Username)); err != nil {
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					if err := os.RemoveAll(fmt.Sprintf("/home/%s/.malino/modules", currentUser.Username)); err != nil {
						return err
					}
				}

				fifthLink, _ := doc.Find("a").Eq(4).Attr("href") // get the fourth link
				fifthLinkURL := linkURL + fifthLink

				if err := downloadFile(fifthLinkURL, "mods.deb"); err != nil {
					return err
				}

				if err := extractWith7z("mods.deb"); err != nil {
					return err
				}

				if err := extractWith7z("data.tar"); err != nil {
					// some symlink bs makes 7zip always return 2.
					if !strings.Contains(err.Error(), "2") {
						return err
					}
				}

				// we have a perfectly good move binary, why not use it :sob:
				if err := execCmd(true, "/bin/bash", "-c", "mv lib/modules/*/kernel /home/$(whoami)/.malino/modules"); err != nil {
					return err
				}

				// go error handling moment
				if err := os.RemoveAll("boot"); err != nil {
					return err
				}

				if err := os.RemoveAll("usr"); err != nil {
					return err
				}

				if err := os.RemoveAll("lib"); err != nil {
					return err
				}

				if err := os.Remove("control.tar"); err != nil {
					return err
				}

				if err := os.Remove("data.tar"); err != nil {
					return err
				}

				if err := os.Remove("mods.deb"); err != nil {
					return err
				}

				// unextract 'em all
				fmt.Printf("Linux kernel module pack \"%s\" downloaded, Now unextracting modules... (this will take a very long time)\n", doc.Find("a").Eq(4).Text())

				zstFiles, err := findZstFiles(fmt.Sprintf("/home/%s/.malino/modules", currentUser.Username))
				if err != nil {
					return fmt.Errorf("error finding compressed files: %s", err.Error())
				}

				err = decompressZstFiles(zstFiles)
				if err != nil {
					return fmt.Errorf("error decompressing compressed files: %s", err.Error())
				}

				err = removeZstFiles(zstFiles)
				if err != nil {
					return fmt.Errorf("error removing compressed files: %s", err.Error())
				}
			}

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

func findZstFiles(root string) ([]string, error) {
	var zstFiles []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zst") {
			zstFiles = append(zstFiles, path)
		}
		return nil
	})
	return zstFiles, err
}

func decompressZstFiles(zstFiles []string) error {
	for _, zstFile := range zstFiles {
		err := execCmd(true, "unzstd", zstFile)
		if err != nil {
			return fmt.Errorf("error decompressing file %s: %w", zstFile, err)
		}
		fmt.Print("\033[1A")
	}
	return nil
}

func removeZstFiles(zstFiles []string) error {
	for _, zstFile := range zstFiles {
		err := os.Remove(zstFile)
		if err != nil {
			return fmt.Errorf("error removing file %s: %w", zstFile, err)
		}
	}
	return nil
}
