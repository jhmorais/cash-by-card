package partner

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	repoClient "github.com/jhmorais/cash-by-card/internal/repositories/client"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

// ---- mocks próprios deste use case: embed da interface + campos Func ----

type mockClientRepoPartnerClients struct {
	repoClient.ClientRepository
	findByPartnerID func(ctx context.Context, partnerID int, name string) ([]*entities.Client, error)
	findByID        func(ctx context.Context, id int) (*entities.Client, error)
	create          func(ctx context.Context, entity *entities.Client) error
	update          func(ctx context.Context, entity *entities.Client) error
}

func (m *mockClientRepoPartnerClients) FindClientByPartnerID(ctx context.Context, partnerID int, name string) ([]*entities.Client, error) {
	if m.findByPartnerID == nil {
		return nil, nil
	}
	return m.findByPartnerID(ctx, partnerID, name)
}

func (m *mockClientRepoPartnerClients) FindClientByID(ctx context.Context, id int) (*entities.Client, error) {
	if m.findByID == nil {
		return nil, nil
	}
	return m.findByID(ctx, id)
}

func (m *mockClientRepoPartnerClients) CreateClient(ctx context.Context, entity *entities.Client) error {
	if m.create == nil {
		return nil
	}
	return m.create(ctx, entity)
}

func (m *mockClientRepoPartnerClients) UpdateClient(ctx context.Context, entity *entities.Client) error {
	if m.update == nil {
		return nil
	}
	return m.update(ctx, entity)
}

type mockPartnerRepoPartnerClients struct {
	repoPartner.PartnerRepository
	findByEmail func(ctx context.Context, email string) ([]*entities.Partner, error)
}

func (m *mockPartnerRepoPartnerClients) FindPartnerByEmail(ctx context.Context, email string) ([]*entities.Partner, error) {
	if m.findByEmail == nil {
		return nil, nil
	}
	return m.findByEmail(ctx, email)
}

// partnerRepoByEmail simula um user com entidade parceira vinculada (gorm Find: slice com 1).
func partnerRepoByEmail(id int) *mockPartnerRepoPartnerClients {
	return &mockPartnerRepoPartnerClients{
		findByEmail: func(ctx context.Context, email string) ([]*entities.Partner, error) {
			return []*entities.Partner{{ID: id, Email: email}}, nil
		},
	}
}

// partnerRepoSemVinculo simula user sem entidade parceira (gorm Find: slice vazio, sem erro).
func partnerRepoSemVinculo() *mockPartnerRepoPartnerClients {
	return &mockPartnerRepoPartnerClients{
		findByEmail: func(ctx context.Context, email string) ([]*entities.Partner, error) {
			return []*entities.Partner{}, nil
		},
	}
}

func TestListClients_Escopo(t *testing.T) {
	t.Run("busca clientes pelo ID da entidade do parceiro resolvido pelo email", func(t *testing.T) {
		var gotPartnerID int
		clientRepo := &mockClientRepoPartnerClients{
			findByPartnerID: func(ctx context.Context, partnerID int, name string) ([]*entities.Client, error) {
				gotPartnerID = partnerID
				return []*entities.Client{{ID: 1, Name: "João"}}, nil
			},
		}
		uc := NewPartnerClientsUseCase(clientRepo, partnerRepoByEmail(7))

		out, err := uc.ListClients(context.Background(), "user@x.com")
		if err != nil {
			t.Fatalf("expected no error, got '%v'", err)
		}
		if gotPartnerID != 7 {
			t.Fatalf("esperado FindClientByPartnerID com partnerID=7, got %d", gotPartnerID)
		}
		if out == nil || len(out.Clients) != 1 || out.Clients[0].Name != "João" {
			t.Fatalf("esperado 1 cliente 'João', got %+v", out)
		}
	})

	t.Run("email sem entidade parceira devolve lista vazia e não consulta clientes", func(t *testing.T) {
		clientRepoCalled := false
		clientRepo := &mockClientRepoPartnerClients{
			findByPartnerID: func(ctx context.Context, partnerID int, name string) ([]*entities.Client, error) {
				clientRepoCalled = true
				return nil, nil
			},
		}
		uc := NewPartnerClientsUseCase(clientRepo, partnerRepoSemVinculo())

		out, err := uc.ListClients(context.Background(), "sem-parceiro@x.com")
		if err != nil {
			t.Fatalf("expected no error, got '%v'", err)
		}
		if clientRepoCalled {
			t.Fatal("repo de clientes não deveria ser consultado sem entidade parceira")
		}
		if out == nil {
			t.Fatal("esperado output não nulo com lista vazia")
		}
		if len(out.Clients) != 0 {
			t.Fatalf("esperado lista vazia, got %+v", out.Clients)
		}
	})
}

func TestCreateClient_VinculoAutomatico(t *testing.T) {
	t.Run("cliente nasce vinculado à entidade do parceiro logado", func(t *testing.T) {
		var captured *entities.Client
		clientRepo := &mockClientRepoPartnerClients{
			create: func(ctx context.Context, entity *entities.Client) error {
				entity.ID = 42 // ID gerado pelo banco
				captured = entity
				return nil
			},
		}
		uc := NewPartnerClientsUseCase(clientRepo, partnerRepoByEmail(7))

		out, err := uc.CreateClient(context.Background(), "user@x.com", &input.CreateClient{
			Name:      "Maria",
			PixType:   1,
			PixKey:    "chave-maria",
			Phone:     "1199999999",
			CPF:       "12345678900",
			Documents: "rg",
		})
		if err != nil {
			t.Fatalf("expected no error, got '%v'", err)
		}
		if captured == nil {
			t.Fatal("CreateClient do repositório não foi chamado")
		}
		if captured.PartnerID == nil || *captured.PartnerID != 7 {
			t.Fatalf("esperado PartnerID=7 (entidade do parceiro logado), got %v", captured.PartnerID)
		}
		if captured.Name != "Maria" || captured.CPF != "12345678900" {
			t.Fatalf("campos do input não repassados: %+v", captured)
		}
		if out == nil || out.ClientID != 42 {
			t.Fatalf("esperado ClientID=42 (ID atribuído pelo mock), got %+v", out)
		}
	})

	t.Run("email sem entidade parceira não cria cliente", func(t *testing.T) {
		createCalled := false
		clientRepo := &mockClientRepoPartnerClients{
			create: func(ctx context.Context, entity *entities.Client) error {
				createCalled = true
				return nil
			},
		}
		uc := NewPartnerClientsUseCase(clientRepo, partnerRepoSemVinculo())

		_, err := uc.CreateClient(context.Background(), "sem-parceiro@x.com", &input.CreateClient{Name: "Maria"})
		if err == nil || !strings.Contains(err.Error(), "nenhum parceiro vinculado") {
			t.Fatalf("esperado erro 'nenhum parceiro vinculado', got '%v'", err)
		}
		if createCalled {
			t.Fatal("repositório de clientes não deveria ser chamado sem entidade parceira")
		}
	})
}

func TestUpdateClient_Regra24h(t *testing.T) {
	partnerID := 7
	outroPartner := 99

	cases := []struct {
		name        string
		partnerRepo *mockPartnerRepoPartnerClients
		client      *entities.Client
		findClient  bool
		wantErr     string
		wantUpdated bool
	}{
		{
			name:        "criado há (janela - 1min) atualiza",
			partnerRepo: partnerRepoByEmail(7),
			client:      &entities.Client{ID: 33, PartnerID: &partnerID, CreatedAt: time.Now().Add(-(partnerEditWindow - time.Minute))},
			findClient:  true,
			wantUpdated: true,
		},
		{
			name:        "criado há (janela + 1min) bloquea com erro 24h",
			partnerRepo: partnerRepoByEmail(7),
			client:      &entities.Client{ID: 33, PartnerID: &partnerID, CreatedAt: time.Now().Add(-(partnerEditWindow + time.Minute))},
			findClient:  true,
			wantErr:     "24h",
		},
		{
			name:        "cliente de outro parceiro nega permissão mesmo dentro da janela",
			partnerRepo: partnerRepoByEmail(7),
			client:      &entities.Client{ID: 33, PartnerID: &outroPartner, CreatedAt: time.Now().Add(-time.Hour)},
			findClient:  true,
			wantErr:     "sem permissão",
		},
		{
			name:        "cliente sem parceiro nega permissão mesmo dentro da janela",
			partnerRepo: partnerRepoByEmail(7),
			client:      &entities.Client{ID: 33, PartnerID: nil, CreatedAt: time.Now().Add(-time.Hour)},
			findClient:  true,
			wantErr:     "sem permissão",
		},
		{
			name:        "cliente inexistente retorna cliente não encontrado",
			partnerRepo: partnerRepoByEmail(7),
			client:      nil,
			findClient:  true,
			wantErr:     "cliente não encontrado",
		},
		{
			name:        "parceiro não resolvido retorna nenhum parceiro vinculado",
			partnerRepo: partnerRepoSemVinculo(),
			client:      &entities.Client{ID: 33, PartnerID: &partnerID, CreatedAt: time.Now().Add(-time.Hour)},
			findClient:  true,
			wantErr:     "nenhum parceiro vinculado",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			updated := false
			clientRepo := &mockClientRepoPartnerClients{
				findByID: func(ctx context.Context, id int) (*entities.Client, error) {
					if !tc.findClient {
						return nil, nil
					}
					return tc.client, nil
				},
				update: func(ctx context.Context, entity *entities.Client) error {
					updated = true
					return nil
				},
			}
			uc := NewPartnerClientsUseCase(clientRepo, tc.partnerRepo)

			_, err := uc.UpdateClient(context.Background(), "user@x.com", &input.UpdateClient{ID: 33, Name: "Novo Nome"})
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got '%v'", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("esperado erro contendo '%s', got '%v'", tc.wantErr, err)
			}
			if updated != tc.wantUpdated {
				t.Fatalf("UpdateClient chamado=%v, esperado %v", updated, tc.wantUpdated)
			}
		})
	}
}

func TestUpdateClient_Sucesso(t *testing.T) {
	partnerID := 7
	existing := &entities.Client{
		ID:        33,
		Name:      "Nome Antigo",
		PixType:   1,
		PixKey:    "chave-antiga",
		Phone:     "1111111111",
		CPF:       "00000000000",
		Documents: "rg-antigo",
		PartnerID: &partnerID,
		CreatedAt: time.Now().Add(-time.Hour), // dentro da janela de 24h
	}
	var saved *entities.Client
	clientRepo := &mockClientRepoPartnerClients{
		findByID: func(ctx context.Context, id int) (*entities.Client, error) {
			return existing, nil
		},
		update: func(ctx context.Context, entity *entities.Client) error {
			saved = entity
			return nil
		},
	}
	uc := NewPartnerClientsUseCase(clientRepo, partnerRepoByEmail(7))

	out, err := uc.UpdateClient(context.Background(), "user@x.com", &input.UpdateClient{
		ID:        33,
		Name:      "Nome Novo",
		PixType:   2,
		PixKey:    "chave-nova",
		Phone:     "2222222222",
		CPF:       "11111111111",
		Documents: "rg-novo",
	})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if saved == nil {
		t.Fatal("UpdateClient do repositório não foi chamado")
	}
	if saved.Name != "Nome Novo" {
		t.Fatalf("esperado Name 'Nome Novo' persistido, got '%s'", saved.Name)
	}
	if saved.PixKey != "chave-nova" || saved.Phone != "2222222222" || saved.CPF != "11111111111" || saved.Documents != "rg-novo" {
		t.Fatalf("campos não atualizados: %+v", saved)
	}
	if out == nil || out.ClientID != 33 {
		t.Fatalf("esperado ClientID=33, got %+v", out)
	}
}
