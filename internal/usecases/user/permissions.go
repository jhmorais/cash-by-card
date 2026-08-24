package user

import (
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

const (
	RoleOrganization = "organization"
	RoleAdmin        = "admin"
	RolePartner      = "partner"
)

// canListUsers: organization e admin acessam a administração de usuários.
func canListUsers(requester *entities.User) error {
	if requester == nil || requester.Email == "" {
		return fmt.Errorf("sem permissão para acessar a administração de usuários")
	}
	if requester.Role != RoleOrganization && requester.Role != RoleAdmin {
		return fmt.Errorf("sem permissão para acessar a administração de usuários")
	}
	return nil
}

// canManageTarget: organization gerencia admin e partner; admin gerencia apenas partner.
// Ninguém gerencia a si mesmo via ações administrativas.
func canManageTarget(requester *entities.User, target *entities.User) error {
	if err := canListUsers(requester); err != nil {
		return err
	}
	if requester.ID == target.ID {
		return fmt.Errorf("não é possível executar esta ação sobre o próprio usuário")
	}
	if requester.Role == RoleAdmin && target.Role != RolePartner {
		return fmt.Errorf("admin só pode gerenciar usuários partner")
	}
	return nil
}

// canAssignRole: organization cria/edita admin e partner; admin apenas partner.
func canAssignRole(requester *entities.User, role string) error {
	if err := canListUsers(requester); err != nil {
		return err
	}
	if role != RoleAdmin && role != RolePartner {
		return fmt.Errorf("role inválida para cadastro: use admin ou partner")
	}
	if requester.Role == RoleAdmin && role != RolePartner {
		return fmt.Errorf("admin só pode cadastrar usuários partner")
	}
	return nil
}
