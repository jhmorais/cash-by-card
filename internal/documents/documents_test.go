package documents

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// magic headers reais dos formatos aceitos
var (
	pdfHead = []byte("%PDF-1.4\n%\xe2\xe3\xcf\xd3")
	jpgHead = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	pngHead = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	gifHead = []byte("GIF89a")
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return &Store{Dir: t.TempDir()}
}

func TestDetectType(t *testing.T) {
	cases := []struct {
		name  string
		head  []byte
		wants string
	}{
		{"pdf", pdfHead, "pdf"},
		{"jpeg", jpgHead, "jpeg"},
		{"png", pngHead, "png"},
		{"gif rejeitado", gifHead, ""},
		{"txt rejeitado", []byte("hello world"), ""},
		{"heic rejeitado", []byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63}, ""},
		{"vazio", nil, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DetectType(tc.head); got != tc.wants {
				t.Fatalf("DetectType(%v) = %q, wants %q", tc.head, got, tc.wants)
			}
		})
	}
}

func TestSaveWritesFileWithDetectedExtension(t *testing.T) {
	store := newTestStore(t)
	content := append(append([]byte{}, pdfHead...), []byte("resto do pdf")...)

	// head é o prefixo já consumido do stream; r continua depois dele
	buf := bytes.NewReader(content)
	head := make([]byte, len(pdfHead))
	if _, err := io.ReadFull(buf, head); err != nil {
		t.Fatalf("setup: %v", err)
	}

	ext, err := store.Save("123.456.789-09", head, buf, int64(len(content)))
	if err != nil {
		t.Fatalf("Save falhou: %v", err)
	}
	if ext != "pdf" {
		t.Fatalf("ext = %q, wants %q", ext, "pdf")
	}

	path := store.Path("123.456.789-09", "pdf")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("arquivo não foi salvo: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("conteúdo gravado difere: %d bytes", len(got))
	}
}

func TestSaveReplacesAndRemovesOldExtension(t *testing.T) {
	store := newTestStore(t)

	pngBody := append([]byte{}, pngHead...)
	if _, err := store.Save("111.222.333-44", nil, bytes.NewReader(pngBody), int64(len(pngBody))); err != nil {
		t.Fatalf("Save png falhou: %v", err)
	}

	pdfBody := append(append([]byte{}, pdfHead...), []byte("x")...)
	ext, err := store.Save("111.222.333-44", pdfHead, bytes.NewReader(pdfBody[14:]), int64(len(pdfBody)))
	if err != nil {
		t.Fatalf("Save pdf falhou: %v", err)
	}
	if ext != "pdf" {
		t.Fatalf("ext = %q, wants %q", ext, "pdf")
	}

	if _, err := os.Stat(store.Path("111.222.333-44", "png")); !os.IsNotExist(err) {
		t.Fatal("arquivo antigo .png deveria ter sido removido")
	}
	if _, err := os.Stat(store.Path("111.222.333-44", "pdf")); err != nil {
		t.Fatalf("arquivo novo .pdf deveria existir: %v", err)
	}

	// apenas um arquivo do cliente no diretório
	matches, _ := filepath.Glob(filepath.Join(store.Dir, "111.222.333-44.*"))
	if len(matches) != 1 {
		t.Fatalf("esperava 1 arquivo, achei %d: %v", len(matches), matches)
	}
}

func TestSaveRejectsUnsupportedFormat(t *testing.T) {
	store := newTestStore(t)
	body := append([]byte{}, gifHead...)

	_, err := store.Save("123", gifHead, bytes.NewReader(body), int64(len(body)))
	if err == nil {
		t.Fatal("esperava erro para gif")
	}
	if !strings.Contains(err.Error(), "formato não suportado") {
		t.Fatalf("mensagem inesperada: %v", err)
	}
	if !errors.Is(err, ErrUnsupportedFormat) {
		t.Fatal("erro deveria ser ErrUnsupportedFormat")
	}
}

func TestSaveRejectsInvalidCPF(t *testing.T) {
	store := newTestStore(t)

	for _, cpf := range []string{"", "../escape", "a/b"} {
		if _, err := store.Save(cpf, pdfHead, bytes.NewReader(pdfHead), int64(len(pdfHead))); err == nil {
			t.Fatalf("esperava erro para cpf %q", cpf)
		}
	}
}

func TestDeleteRemovesAllFilesAndToleratesMissing(t *testing.T) {
	store := newTestStore(t)
	body := append([]byte{}, jpgHead...)
	if _, err := store.Save("999.888.777-66", nil, bytes.NewReader(body), int64(len(body))); err != nil {
		t.Fatalf("Save falhou: %v", err)
	}

	if err := store.Delete("999.888.777-66"); err != nil {
		t.Fatalf("Delete falhou: %v", err)
	}
	if _, err := os.Stat(store.Path("999.888.777-66", "jpeg")); !os.IsNotExist(err) {
		t.Fatal("arquivo deveria ter sido apagado")
	}

	// ausente não é erro
	if err := store.Delete("999.888.777-66"); err != nil {
		t.Fatalf("Delete de arquivo inexistente não deveria errar: %v", err)
	}
	if err := store.Delete("000.000.000-00"); err != nil {
		t.Fatalf("Delete de cpf sem documento não deveria errar: %v", err)
	}
}

func TestNewStoreUsesEnvAndCreatesDir(t *testing.T) {
	t.Setenv("DOCS_DIR", filepath.Join(t.TempDir(), "nested", "docs"))

	store := NewStore()
	if store.Dir != os.Getenv("DOCS_DIR") {
		t.Fatalf("Dir = %q, wants %q", store.Dir, os.Getenv("DOCS_DIR"))
	}
	if info, err := os.Stat(store.Dir); err != nil || !info.IsDir() {
		t.Fatalf("diretório deveria ter sido criado: %v", err)
	}
}

func TestNewStoreDefaultsToDocsDir(t *testing.T) {
	t.Setenv("DOCS_DIR", "")

	store := NewStore()
	if store.Dir != "./docs" {
		t.Fatalf("Dir default = %q, wants ./docs", store.Dir)
	}
}

func TestSaveIsSeekFriendly(t *testing.T) {
	store := newTestStore(t)
	// leitor que só entrega 1 byte por vez (io.Reader custom, não seeker)
	// e head vazia: o conteúdo inteiro vem pelo reader
	chunked := iotestReader(pngHead)
	if _, err := store.Save("555.555.555-55", nil, chunked, int64(len(pngHead))); err != nil {
		t.Fatalf("Save falhou: %v", err)
	}
	got, err := os.ReadFile(store.Path("555.555.555-55", "png"))
	if err != nil {
		t.Fatalf("arquivo não salvo: %v", err)
	}
	if !bytes.Equal(got, pngHead) {
		t.Fatalf("conteúdo difere")
	}
}

func iotestReader(b []byte) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		for _, c := range b {
			pw.Write([]byte{c})
		}
		pw.Close()
	}()
	return pr
}
