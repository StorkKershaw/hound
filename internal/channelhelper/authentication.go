package channelhelper

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/StorkKershaw/hound/internal/browser"
	"github.com/mxschmitt/playwright-go"
)

const authRetryInterval = 30 * time.Second

func Authenticate(in <-chan struct{}, out chan<- error) {
	for {
		err := browser.Use("chrome", func(page playwright.Page) {
			for range in {
				if err := authenticate(page); err != nil {
					out <- err
				}
			}
		})

		if err == nil {
			return
		}

		log.Printf("browser session failed, retrying in %s: %v", authRetryInterval, err)
		out <- err
		time.Sleep(authRetryInterval)
	}
}

func authenticate(page playwright.Page) error {
	response, err := page.Goto("http://www.msftconnecttest.com/connecttest.txt")
	if err != nil {
		return err
	}
	if response == nil {
		return fmt.Errorf("no response from connectivity check")
	}

	content, err := response.Text()
	if err != nil {
		return err
	}

	if content == "Microsoft Connect Test" {
		return nil
	}

	parsed, err := url.Parse(page.URL())
	if err != nil {
		return err
	}

	if hostname := parsed.Hostname(); hostname != "service.wi2.ne.jp" {
		return fmt.Errorf("unexpected hostname: %s", hostname)
	}

	locator := page.GetByRole(
		*playwright.AriaRoleButton,
		playwright.PageGetByRoleOptions{Name: "同意する"},
	)

	if err := locator.WaitFor(); err != nil {
		return err
	}

	return locator.Click()
}
