-- 2026-08-26: documento único por cliente (PDF/JPEG/PNG)
-- o arquivo fica no filesystem (DOCS_DIR, default ./docs) como {cpf}.{ext};
-- aqui ficam só os metadados. document_type guarda a extensão detectada
-- por magic bytes (pdf/jpeg/png); NULL/'' = cliente sem documento.
ALTER TABLE `client`
  ADD COLUMN `document_name` VARCHAR(255) NULL,
  ADD COLUMN `document_type` VARCHAR(50) NULL,
  ADD COLUMN `document_size` INT NULL;
