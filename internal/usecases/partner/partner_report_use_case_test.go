package partner

import (
	"context"
	"testing"
	"time"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	repoLoan "github.com/jhmorais/cash-by-card/internal/repositories/loan"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

// ---- mocks próprios deste use case: embed da interface + campos Func ----

type mockLoanRepoPartnerReport struct {
	repoLoan.LoanRepository
	findByPartnerID func(ctx context.Context, partnerID int) ([]*entities.Loan, error)
}

func (m *mockLoanRepoPartnerReport) FindLoanByPartnerID(ctx context.Context, partnerID int) ([]*entities.Loan, error) {
	if m.findByPartnerID == nil {
		return nil, nil
	}
	return m.findByPartnerID(ctx, partnerID)
}

type mockPartnerRepoPartnerReport struct {
	repoPartner.PartnerRepository
	findByEmail func(ctx context.Context, email string) ([]*entities.Partner, error)
}

func (m *mockPartnerRepoPartnerReport) FindPartnerByEmail(ctx context.Context, email string) ([]*entities.Partner, error) {
	if m.findByEmail == nil {
		return nil, nil
	}
	return m.findByEmail(ctx, email)
}

func TestPartnerReport_Agregacao(t *testing.T) {
	now := time.Now()
	currentMonth := int(now.Month())

	// 2 empréstimos no mês corrente, 1 no ano passado, 1 dois anos atrás.
	loans := []*entities.Loan{
		{ID: 101, PartnerAmount: 100.50, CreatedAt: time.Date(now.Year(), now.Month(), 3, 10, 0, 0, 0, time.UTC), Client: entities.Client{Name: "Maria Silva"}},
		{ID: 102, PartnerAmount: 50, CreatedAt: time.Date(now.Year(), now.Month(), 15, 14, 30, 0, 0, time.UTC), Client: entities.Client{Name: ""}},
		{ID: 103, PartnerAmount: 200, CreatedAt: time.Date(now.Year()-1, 6, 10, 9, 0, 0, 0, time.UTC), Client: entities.Client{Name: "João Souza"}},
		{ID: 104, PartnerAmount: 75.25, CreatedAt: time.Date(now.Year()-2, 1, 20, 8, 0, 0, 0, time.UTC), Client: entities.Client{Name: "Ana Lima"}},
	}

	var gotPartnerID int
	loanRepo := &mockLoanRepoPartnerReport{
		findByPartnerID: func(ctx context.Context, partnerID int) ([]*entities.Loan, error) {
			gotPartnerID = partnerID
			return loans, nil
		},
	}
	partnerRepo := &mockPartnerRepoPartnerReport{
		findByEmail: func(ctx context.Context, email string) ([]*entities.Partner, error) {
			return []*entities.Partner{{ID: 7, Email: email}}, nil
		},
	}
	uc := NewPartnerReportUseCase(loanRepo, partnerRepo)

	report, err := uc.Execute(context.Background(), "parceiro@x.com")
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if report == nil {
		t.Fatal("esperado report não nulo")
	}
	if gotPartnerID != 7 {
		t.Fatalf("esperado FindLoanByPartnerID com partnerID=7, got %d", gotPartnerID)
	}

	t.Run("summary soma todos os empréstimos", func(t *testing.T) {
		if report.Summary.TotalLoans != 4 {
			t.Fatalf("esperado TotalLoans=4, got %d", report.Summary.TotalLoans)
		}
		if report.Summary.TotalCommission != 425.75 {
			t.Fatalf("esperado TotalCommission=425.75, got %v", report.Summary.TotalCommission)
		}
	})

	t.Run("annual vai do ano mais antigo ao corrente com 12 meses cada", func(t *testing.T) {
		if len(report.Annual) != 3 {
			t.Fatalf("esperado 3 anos (de %d até %d), got %d: %+v", now.Year()-2, now.Year(), len(report.Annual), report.Annual)
		}
		for i, wantYear := range []int{now.Year() - 2, now.Year() - 1, now.Year()} {
			py := report.Annual[i]
			if py.Year != wantYear {
				t.Fatalf("ano na posição %d: esperado %d, got %d", i, wantYear, py.Year)
			}
			if len(py.Months) != 12 {
				t.Fatalf("ano %d: esperado 12 meses, got %d", wantYear, len(py.Months))
			}
			for m := 1; m <= 12; m++ {
				if py.Months[m-1].Month != m {
					t.Fatalf("ano %d: mês na posição %d deveria ser %d, got %+v", wantYear, m-1, m, py.Months[m-1])
				}
			}
		}
	})

	t.Run("meses sem empréstimos ficam zerados", func(t *testing.T) {
		zero := report.Annual[0].Months[11] // dezembro de dois anos atrás
		if zero.Loans != 0 || zero.Commission != 0 {
			t.Fatalf("esperado mês 12/%d zerado, got %+v", now.Year()-2, zero)
		}
	})

	t.Run("mês corrente agrega os 2 empréstimos e a comissão somada", func(t *testing.T) {
		m := report.Annual[2].Months[currentMonth-1]
		if m.Loans != 2 {
			t.Fatalf("esperado Loans=2 no mês corrente, got %+v", m)
		}
		if m.Commission != 150.50 {
			t.Fatalf("esperado Commission=150.50 no mês corrente, got %v", m.Commission)
		}
	})

	t.Run("meses de anos anteriores agregam seus empréstimos", func(t *testing.T) {
		lastYear := report.Annual[1].Months[5] // junho do ano passado
		if lastYear.Loans != 1 || lastYear.Commission != 200 {
			t.Fatalf("esperado junho/%d com 1 empréstimo e comissão 200, got %+v", now.Year()-1, lastYear)
		}
		oldest := report.Annual[0].Months[0] // janeiro de dois anos atrás
		if oldest.Loans != 1 || oldest.Commission != 75.25 {
			t.Fatalf("esperado janeiro/%d com 1 empréstimo e comissão 75.25, got %+v", now.Year()-2, oldest)
		}
	})

	t.Run("currentMonth traz exatamente os 2 empréstimos do mês corrente", func(t *testing.T) {
		if len(report.CurrentMonth) != 2 {
			t.Fatalf("esperado 2 empréstimos no mês corrente, got %d: %+v", len(report.CurrentMonth), report.CurrentMonth)
		}
		type detail struct {
			commission float64
			clientName string
		}
		byLoanID := map[int]detail{}
		for _, d := range report.CurrentMonth {
			byLoanID[d.LoanID] = detail{commission: d.Commission, clientName: d.ClientName}
		}
		if len(byLoanID) != 2 {
			t.Fatalf("esperado 2 LoanIDs distintos, got %v", byLoanID)
		}
		d101, ok := byLoanID[101]
		if !ok || d101.clientName != "Maria Silva" || d101.commission != 100.50 {
			t.Fatalf("empréstimo 101 com dados errados: %+v", d101)
		}
		d102, ok := byLoanID[102]
		if !ok || d102.clientName != "" || d102.commission != 50 {
			t.Fatalf("empréstimo 102 com dados errados: %+v", d102)
		}
		if _, ok := byLoanID[103]; ok {
			t.Fatal("empréstimo de outro ano não deveria aparecer no mês corrente")
		}
	})

	t.Run("generatedAt preenchido", func(t *testing.T) {
		if report.GeneratedAt.IsZero() {
			t.Fatal("esperado GeneratedAt preenchido")
		}
	})
}

func TestPartnerReport_SemEntidadeParceira(t *testing.T) {
	loanRepoCalled := false
	loanRepo := &mockLoanRepoPartnerReport{
		findByPartnerID: func(ctx context.Context, partnerID int) ([]*entities.Loan, error) {
			loanRepoCalled = true
			return nil, nil
		},
	}
	partnerRepo := &mockPartnerRepoPartnerReport{
		findByEmail: func(ctx context.Context, email string) ([]*entities.Partner, error) {
			return []*entities.Partner{}, nil // gorm Find: slice vazio, sem erro
		},
	}
	uc := NewPartnerReportUseCase(loanRepo, partnerRepo)

	report, err := uc.Execute(context.Background(), "sem-parceiro@x.com")
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if report == nil {
		t.Fatal("esperado report não nulo (vazio)")
	}
	if loanRepoCalled {
		t.Fatal("repo de empréstimos não deveria ser consultado sem entidade parceira")
	}
	if report.Summary.TotalLoans != 0 || report.Summary.TotalCommission != 0 {
		t.Fatalf("esperado Summary zerado, got %+v", report.Summary)
	}
	if len(report.Annual) != 0 {
		t.Fatalf("esperado Annual vazio, got %d anos", len(report.Annual))
	}
	if len(report.CurrentMonth) != 0 {
		t.Fatalf("esperado CurrentMonth vazio, got %d itens", len(report.CurrentMonth))
	}
}

func TestPartnerReport_PropagaErros(t *testing.T) {
	t.Run("erro ao buscar parceiro", func(t *testing.T) {
		partnerRepo := &mockPartnerRepoPartnerReport{
			findByEmail: func(ctx context.Context, email string) ([]*entities.Partner, error) {
				return nil, context.DeadlineExceeded
			},
		}
		uc := NewPartnerReportUseCase(&mockLoanRepoPartnerReport{}, partnerRepo)
		report, err := uc.Execute(context.Background(), "parceiro@x.com")
		if err == nil {
			t.Fatal("esperado erro propagado de FindPartnerByEmail")
		}
		if report != nil {
			t.Fatalf("esperado report nil junto com erro, got %+v", report)
		}
	})

	t.Run("erro ao buscar empréstimos", func(t *testing.T) {
		loanRepo := &mockLoanRepoPartnerReport{
			findByPartnerID: func(ctx context.Context, partnerID int) ([]*entities.Loan, error) {
				return nil, context.DeadlineExceeded
			},
		}
		partnerRepo := &mockPartnerRepoPartnerReport{
			findByEmail: func(ctx context.Context, email string) ([]*entities.Partner, error) {
				return []*entities.Partner{{ID: 7, Email: email}}, nil
			},
		}
		uc := NewPartnerReportUseCase(loanRepo, partnerRepo)
		report, err := uc.Execute(context.Background(), "parceiro@x.com")
		if err == nil {
			t.Fatal("esperado erro propagado de FindLoanByPartnerID")
		}
		if report != nil {
			t.Fatalf("esperado report nil junto com erro, got %+v", report)
		}
	})
}
