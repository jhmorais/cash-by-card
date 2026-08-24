-- 2026-08-24: primeiro acesso e administracao de usuarios
-- password NULL = usuario pendente de primeiro acesso
--
-- BOOTSTRAP OBRIGATORIO EM PRODUCAO (apos rodar este script):
-- ninguem pode criar/atribuir o role 'organization' pela API (por design),
-- entao promova manualmente a conta da organizacao ANTES do primeiro uso:
--   UPDATE `user` SET role='organization' WHERE email='<email-do-dono>';
ALTER TABLE `user` MODIFY `password` VARCHAR(100) NULL;

CREATE TABLE `password_reset_token` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `token_hash` VARCHAR(64) NOT NULL,
  `expires_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `used_at` TIMESTAMP NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `password_reset_token_user_id` (`user_id`),
  KEY `password_reset_token_token_hash` (`token_hash`),
  CONSTRAINT `password_reset_token_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
