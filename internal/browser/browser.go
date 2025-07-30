package browser

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func Use(channel string, fn func(page playwright.Page)) {
	options := playwright.RunOptions{
		SkipInstallBrowsers: true,
	}

	if err := playwright.Install(&options); err != nil {
		log.Printf("could not install playwright dependencies: %v", err)
		return
	}

	pw, err := playwright.Run(&options)
	if err != nil {
		log.Printf("could not start playwright: %v", err)
		return
	}
	defer pw.Stop()

	context, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Channel: playwright.String(channel),
		},
	)
	if err != nil {
		log.Printf("could not launch browser: %v", err)
		return
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		log.Printf("could not create new page: %v", err)
		return
	}
	defer page.Close()

	fn(page)
}
