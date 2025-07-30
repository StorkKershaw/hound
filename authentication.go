package main

import (
	"fmt"
	"net/url"

	"github.com/StorkKershaw/hound/internal/browser"
	"github.com/playwright-community/playwright-go"
)

func authenticate(in <-chan struct{}, out chan<- error) {
	browser.Use("chrome", func(page playwright.Page) {
		for range in {
			response, err := page.Goto("http://www.msftconnecttest.com/connecttest.txt")
			if err != nil {
				out <- err
				continue
			}

			content, err := response.Text()
			if err != nil {
				out <- err
				continue
			}

			if content == "Microsoft Connect Test" {
				continue
			}

			parsed, err := url.Parse(page.URL())
			if err != nil {
				out <- err
				continue
			}

			if hostname := parsed.Hostname(); hostname != "service.wi2.ne.jp" {
				out <- fmt.Errorf("unexpected hostname: %s", hostname)
				continue
			}

			locator := page.GetByRole(
				*playwright.AriaRoleButton,
				playwright.PageGetByRoleOptions{Name: "同意する"},
			)

			if err := locator.WaitFor(); err != nil {
				out <- err
				continue
			}

			if err := locator.Click(); err != nil {
				out <- err
				continue
			}
		}
	})
}
