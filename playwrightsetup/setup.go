package playwrightsetup

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

// PlaywrightManager estrutura que encapsula o Playwright e o navegador
type PlaywrightManager struct {
	Playwright *playwright.Playwright
	Browser    playwright.Browser
}

// NewPlaywrightManager inicializa o Playwright e o navegador
func NewPlaywrightManager() (*PlaywrightManager, error) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Failed to start Playwright: %v", err)
		return nil, err
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		pw.Stop()
		log.Fatalf("Failed to launch Chromium: %v", err)
		return nil, err
	}

	return &PlaywrightManager{
		Playwright: pw,
		Browser:    browser,
	}, nil
}

// Close finaliza o navegador e o Playwright
func (pm *PlaywrightManager) Close() {
	pm.Browser.Close()
	pm.Playwright.Stop()
}
