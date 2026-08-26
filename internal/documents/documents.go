// Package documents armazena o documento único de cada cliente no
// filesystem (DOCS_DIR, default ./docs) como {cpf}.{ext}. É um pacote
// puro de fs — sem gorm — para ser testado isoladamente.
package documents

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ErrUnsupportedFormat é devolvido quando o conteúdo não é PDF, JPEG ou PNG.
var ErrUnsupportedFormat = errors.New("formato não suportado — use PDF, JPEG ou PNG")

// Store grava e recupera documentos em um diretório do filesystem.
type Store struct {
	Dir string
}

// NewStore cria o Store a partir da env DOCS_DIR (default ./docs),
// garantindo que o diretório exista.
func NewStore() *Store {
	dir := os.Getenv("DOCS_DIR")
	if dir == "" {
		dir = "./docs"
	}
	_ = os.MkdirAll(dir, 0o755)
	return &Store{Dir: dir}
}

// DetectType identifica o formato pelos magic bytes do início do arquivo.
// Devolve "pdf", "jpeg" ou "png"; "" significa formato não suportado.
func DetectType(head []byte) string {
	switch {
	case bytes.HasPrefix(head, []byte("%PDF-")):
		return "pdf"
	case len(head) >= 2 && head[0] == 0xFF && head[1] == 0xD8:
		return "jpeg"
	case bytes.HasPrefix(head, []byte{0x89, 'P', 'N', 'G'}):
		return "png"
	default:
		return ""
	}
}

// validCPF aceita apenas os caracteres que aparecem em um CPF
// (dígitos, ponto e traço) para evitar path traversal no nome do arquivo.
func validCPF(cpf string) bool {
	if cpf == "" {
		return false
	}
	for _, c := range cpf {
		if !(c >= '0' && c <= '9' || c == '.' || c == '-') {
			return false
		}
	}
	return true
}

// Save valida o formato pelos magic bytes (head), grava o documento em
// Dir/{cpf}.{ext} de forma atômica (arquivo temporário + rename) e remove
// qualquer {cpf}.* anterior com outra extensão. Devolve a extensão usada.
// head é o prefixo já consumido do stream (r continua logo depois dele);
// se vier vazio, os primeiros bytes são lidos de r.
func (s *Store) Save(cpf string, head []byte, r io.Reader, size int64) (string, error) {
	if !validCPF(cpf) {
		return "", fmt.Errorf("cpf inválido")
	}

	if len(head) == 0 && r != nil {
		buf := make([]byte, 512)
		n, err := io.ReadFull(r, buf)
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return "", fmt.Errorf("falha ao ler documento: %w", err)
		}
		head = buf[:n]
	}

	ext := DetectType(head)
	if ext == "" {
		return "", ErrUnsupportedFormat
	}

	tmp, err := os.CreateTemp(s.Dir, ".upload-*")
	if err != nil {
		return "", fmt.Errorf("falha ao criar arquivo temporário: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op se o rename deu certo

	if _, err := tmp.Write(head); err != nil {
		tmp.Close()
		return "", fmt.Errorf("falha ao gravar documento: %w", err)
	}
	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		return "", fmt.Errorf("falha ao gravar documento: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("falha ao gravar documento: %w", err)
	}

	final := s.Path(cpf, ext)
	if err := os.Rename(tmpName, final); err != nil {
		return "", fmt.Errorf("falha ao mover documento para o destino: %w", err)
	}

	// substituição: apaga {cpf}.* de extensões anteriores
	s.removeOthers(cpf, final)

	return ext, nil
}

// Path devolve o caminho do documento do cliente com a extensão dada.
func (s *Store) Path(cpf string, ext string) string {
	return filepath.Join(s.Dir, cpf+"."+ext)
}

// Delete remove todos os arquivos {cpf}.* do cliente. Arquivo ausente
// não é erro.
func (s *Store) Delete(cpf string) error {
	if !validCPF(cpf) {
		return fmt.Errorf("cpf inválido")
	}

	matches, err := filepath.Glob(filepath.Join(s.Dir, cpf+".*"))
	if err != nil {
		return err
	}
	for _, m := range matches {
		if err := os.Remove(m); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (s *Store) removeOthers(cpf, keep string) {
	matches, err := filepath.Glob(filepath.Join(s.Dir, cpf+".*"))
	if err != nil {
		return
	}
	for _, m := range matches {
		if m == keep {
			continue
		}
		_ = os.Remove(m)
	}
}
