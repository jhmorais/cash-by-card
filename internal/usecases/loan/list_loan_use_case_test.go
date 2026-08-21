package loan

import (
	"context"
	"errors"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

// mockLoanRepository embeds the interface so only ListLoan needs to be implemented.
type mockLoanRepository struct {
	repositories.LoanRepository
	listLoanFunc func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error)
}

func (m *mockLoanRepository) ListLoan(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
	return m.listLoanFunc(ctx, filter, pagination)
}

func TestListLoanUseCase_Execute_DefaultPagination(t *testing.T) {
	var gotPagination *input.Pagination
	var gotFilter *input.ListLoanFilter

	useCase := NewListLoansUseCase(&mockLoanRepository{
		listLoanFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
			gotFilter = filter
			gotPagination = pagination
			return []*entities.Loan{{ID: 1}}, 1, nil
		},
	})

	result, err := useCase.Execute(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if gotPagination == nil || gotPagination.Page != 1 || gotPagination.Limit != 10 {
		t.Fatalf("expected default pagination {Page:1, Limit:10}, got %+v", gotPagination)
	}
	if gotFilter == nil {
		t.Fatal("expected non-nil filter passed to repository")
	}
	if result.Total != 1 || result.Page != 1 || result.Limit != 10 || result.TotalPages != 1 {
		t.Fatalf("unexpected output: %+v", result)
	}
	if len(result.Loans) != 1 {
		t.Fatalf("expected 1 loan, got %d", len(result.Loans))
	}
}

func TestListLoanUseCase_Execute_NormalizesInvalidPagination(t *testing.T) {
	var gotPagination *input.Pagination

	useCase := NewListLoansUseCase(&mockLoanRepository{
		listLoanFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
			gotPagination = pagination
			return []*entities.Loan{}, 0, nil
		},
	})

	_, err := useCase.Execute(context.Background(), nil, &input.Pagination{Page: 0, Limit: 500})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if gotPagination.Page != 1 {
		t.Fatalf("expected page normalized to 1, got %d", gotPagination.Page)
	}
	if gotPagination.Limit != 100 {
		t.Fatalf("expected limit capped to 100, got %d", gotPagination.Limit)
	}
}

func TestListLoanUseCase_Execute_PassesFilterThrough(t *testing.T) {
	var gotFilter *input.ListLoanFilter

	useCase := NewListLoansUseCase(&mockLoanRepository{
		listLoanFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
			gotFilter = filter
			return []*entities.Loan{}, 0, nil
		},
	})

	paymentStatus := "pending"
	clientName := "joão"
	clientCPF := "12345678900"

	filter := &input.ListLoanFilter{
		PaymentStatus: &paymentStatus,
		ClientName:    &clientName,
		ClientCPF:     &clientCPF,
	}

	_, err := useCase.Execute(context.Background(), filter, &input.Pagination{Page: 2, Limit: 25})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if gotFilter != filter {
		t.Fatal("expected the exact same filter instance to be passed to repository")
	}
}

func TestListLoanUseCase_Execute_TotalPages(t *testing.T) {
	useCase := NewListLoansUseCase(&mockLoanRepository{
		listLoanFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
			return []*entities.Loan{}, 42, nil
		},
	})

	result, err := useCase.Execute(context.Background(), nil, &input.Pagination{Page: 3, Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if result.TotalPages != 5 {
		t.Fatalf("expected 5 total pages for 42 items with limit 10, got %d", result.TotalPages)
	}
}

func TestListLoanUseCase_Execute_RepositoryError(t *testing.T) {
	useCase := NewListLoansUseCase(&mockLoanRepository{
		listLoanFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
			return nil, 0, errors.New("db down")
		},
	})

	result, err := useCase.Execute(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}
