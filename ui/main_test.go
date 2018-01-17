package main_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/sclevine/agouti"
)

var sampleHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
})

func TestMain(m *testing.M) {
	go func() {
		log.Fatal(http.ListenAndServe(":8080", sampleHandler))
	}()
	os.Exit(m.Run())
}

func TestHelloWorld(t *testing.T) {
	// driver := agouti.PhantomJS()
	driver := agouti.Selenium()
	// driver = agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		t.Fatal("Failed to start Selenium:", err)
	}
	page, err := agouti.NewPage(driver.URL())
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate("http://localhost:8080"); err != nil {
		t.Fatal("Failed to navigate:", err)
	}

	URL, err := page.URL()
	if err != nil {
		t.Fatal("Failed to get page URL:", err)
	}

	expectedURL := "http://localhost:8080/"
	if URL != expectedURL {
		t.Fatal("Expected URL to be", expectedURL, "but got", URL)
	}

	msg, err := page.Find("#foo").Text()
	if err != nil {
		t.Fatal("Failed to get foo text:", err)
	}

	expectedMsg := "hello world"
	if msg != expectedMsg {
		t.Fatal("Expected hello world to be", expectedMsg, "but got", loginMsg)
	}

	if err := driver.Stop(); err != nil {
		t.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}
