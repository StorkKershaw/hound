package browser

import (
	"fmt"

	pw "github.com/mxschmitt/playwright-go"
)

func Use(channel string, fn func(page pw.Page)) error {
	runOptions := &pw.RunOptions{
		OnlyInstallShell:    true,
		SkipInstallBrowsers: true,
		Browsers:            []string{channel},
	}

	if err := pw.Install(runOptions); err != nil {
		return fmt.Errorf("could not install playwright dependencies: %w", err)
	}

	playwright, err := pw.Run(runOptions)
	if err != nil {
		return fmt.Errorf("could not start playwright: %w", err)
	}
	defer playwright.Stop()

	launchOptions := pw.BrowserTypeLaunchOptions{
		Channel: new(channel),
	}
	context, err := playwright.Chromium.Launch(launchOptions)
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("could not create new page: %w", err)
	}
	defer page.Close()

	fn(page)
	return nil
}
