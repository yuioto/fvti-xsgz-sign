package notify

import (
	"context"
	"net"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestBuildEmailBody(t *testing.T) {
	body := buildEmailBody("me@example.com", []string{"you@example.com"}, "status", "<h1>ok</h1>")
	if body == "" {
		t.Fatal("expected body not empty")
	}
	if !contains(body, "Content-Type: text/html") {
		t.Fatalf("expected HTML content-type, got %q", body)
	}
	if !contains(body, "<h1>ok</h1>") {
		t.Fatalf("expected HTML payload, got %q", body)
	}
}

func TestDialWithProxyDirect(t *testing.T) {
	if err := os.Setenv("ALL_PROXY", ""); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("all_proxy", ""); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("HTTPS_PROXY", ""); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("https_proxy", ""); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("NO_PROXY", ""); err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen error: %v", err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		conn.Close()
	}()

	conn, err := dialWithProxy(context.Background(), "tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("dialWithProxy failed: %v", err)
	}
	conn.Close()
}

func TestGetProxyURL(t *testing.T) {
	_ = os.Setenv("ALL_PROXY", "socks5://127.0.0.1:1080")
	defer os.Unsetenv("ALL_PROXY")

	url := getProxyURL()
	if url == nil {
		t.Fatal("expected proxy URL")
	}
	if url.Scheme != "socks5" {
		t.Fatalf("expected socks5 scheme, got %s", url.Scheme)
	}
}

func TestShouldBypassProxy(t *testing.T) {
	_ = os.Setenv("NO_PROXY", "127.0.0.1")
	defer os.Unsetenv("NO_PROXY")

	u, _ := url.Parse("http://127.0.0.1:8080")
	if !shouldBypassProxy(u, "127.0.0.1:587") {
		t.Fatal("expected bypass proxy for 127.0.0.1")
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
