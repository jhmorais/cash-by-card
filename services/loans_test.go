package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

// mockListLoanUseCase embeds the interface so only Execute needs to be implemented.
type mockListLoanUseCase struct {
	contracts.ListLoanUseCase
	executeFunc func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error)
}

func (m *mockListLoanUseCase) Execute(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
	return m.executeFunc(ctx, filter, pagination)
}

func listLoansHandler(useCase contracts.ListLoanUseCase) *Handler {
	return &Handler{ListLoanUseCase: useCase}
}

func TestListLoans_DefaultParams(t *testing.T) {
	var gotPagination *input.Pagination
	var gotFilter *input.ListLoanFilter

	handler := listLoansHandler(&mockListLoanUseCase{
		executeFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
			gotFilter = filter
			gotPagination = pagination
			return &output.ListLoan{Loans: []*entities.Loan{}, Total: 0, Page: 1, Limit: 10, TotalPages: 0}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/loans", nil)
	rec := httptest.NewRecorder()
	handler.ListLoans(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if gotPagination == nil || gotPagination.Page != 1 || gotPagination.Limit != 10 {
		t.Fatalf("expected default pagination {1, 10}, got %+v", gotPagination)
	}
	if gotFilter == nil {
		t.Fatal("expected non-nil filter")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	for _, key := range []string{"loans", "total", "page", "limit", "totalPages"} {
		if _, ok := body[key]; !ok {
			t.Fatalf("expected key '%s' in response, got %v", key, body)
		}
	}
}

func TestListLoans_ParsesAllFilters(t *testing.T) {
	var gotFilter *input.ListLoanFilter
	var gotPagination *input.Pagination

	handler := listLoansHandler(&mockListLoanUseCase{
		executeFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
			gotFilter = filter
			gotPagination = pagination
			return &output.ListLoan{Loans: []*entities.Loan{}, Total: 0, Page: 2, Limit: 25, TotalPages: 0}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/loans?page=2&limit=25&paymentStatus=pending&type=1&clientName=jo%C3%A3o&partnerName=maria&amountMin=100.5&amountMax=500&askValueMin=10&askValueMax=20&clientAmountMin=1&clientAmountMax=2&grossProfitMin=3&grossProfitMax=4&profitMin=5&profitMax=6&partnerAmountMin=7&partnerAmountMax=8&operationPercentMin=9&operationPercentMax=10&partnerPercentMin=11&partnerPercentMax=12&numberCardsMin=1&numberCardsMax=5", nil)
	rec := httptest.NewRecorder()
	handler.ListLoans(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if gotPagination == nil || gotPagination.Page != 2 || gotPagination.Limit != 25 {
		t.Fatalf("expected pagination {2, 25}, got %+v", gotPagination)
	}

	if gotFilter.PaymentStatus == nil || *gotFilter.PaymentStatus != "pending" {
		t.Fatalf("expected PaymentStatus 'pending', got %+v", gotFilter.PaymentStatus)
	}
	if gotFilter.Type == nil || *gotFilter.Type != 1 {
		t.Fatalf("expected Type 1, got %+v", gotFilter.Type)
	}
	if gotFilter.ClientName == nil || *gotFilter.ClientName != "joão" {
		t.Fatalf("expected ClientName 'joão', got %+v", gotFilter.ClientName)
	}
	if gotFilter.PartnerName == nil || *gotFilter.PartnerName != "maria" {
		t.Fatalf("expected PartnerName 'maria', got %+v", gotFilter.PartnerName)
	}
	if gotFilter.AmountMin == nil || *gotFilter.AmountMin != 100.5 {
		t.Fatalf("expected AmountMin 100.5, got %+v", gotFilter.AmountMin)
	}
	if gotFilter.AmountMax == nil || *gotFilter.AmountMax != 500 {
		t.Fatalf("expected AmountMax 500, got %+v", gotFilter.AmountMax)
	}
	if gotFilter.NumberCardsMin == nil || *gotFilter.NumberCardsMin != 1 {
		t.Fatalf("expected NumberCardsMin 1, got %+v", gotFilter.NumberCardsMin)
	}
	if gotFilter.NumberCardsMax == nil || *gotFilter.NumberCardsMax != 5 {
		t.Fatalf("expected NumberCardsMax 5, got %+v", gotFilter.NumberCardsMax)
	}
	// spot-check the remaining ranges
	if gotFilter.AskValueMin == nil || gotFilter.ProfitMax == nil || gotFilter.PartnerAmountMin == nil ||
		gotFilter.OperationPctMax == nil || gotFilter.PartnerPctMin == nil || gotFilter.GrossProfitMax == nil ||
		gotFilter.ClientAmountMin == nil {
		t.Fatal("expected all numeric range filters to be parsed")
	}
}

func TestListLoans_InvalidParamReturns400(t *testing.T) {
	cases := []struct {
		name string
		url  string
	}{
		{"invalid page", "/admin/loans?page=abc"},
		{"invalid limit", "/admin/loans?limit=x"},
		{"page below 1", "/admin/loans?page=0"},
		{"limit above 100", "/admin/loans?limit=101"},
		{"invalid amountMin", "/admin/loans?amountMin=notanumber"},
		{"invalid type", "/admin/loans?type=3"},
		{"invalid paymentStatus", "/admin/loans?paymentStatus=foo"},
		{"invalid numberCardsMin", "/admin/loans?numberCardsMin=zz"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := listLoansHandler(&mockListLoanUseCase{
				executeFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
					t.Error("use case should not be called for invalid params")
					return nil, nil
				},
			})

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rec := httptest.NewRecorder()
			handler.ListLoans(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d, body: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestListLoans_UseCaseErrorReturns404(t *testing.T) {
	handler := listLoansHandler(&mockListLoanUseCase{
		executeFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
			return nil, errors.New("db down")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/loans", nil)
	rec := httptest.NewRecorder()
	handler.ListLoans(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestListLoans_ParsesAndNormalizesCpfParams(t *testing.T) {
	var gotFilter *input.ListLoanFilter

	handler := listLoansHandler(&mockListLoanUseCase{
		executeFunc: func(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
			gotFilter = filter
			return &output.ListLoan{Loans: []*entities.Loan{}, Total: 0, Page: 1, Limit: 10, TotalPages: 0}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/loans?clientCpf=123.456.789-00&partnerCpf=98765432100", nil)
	rec := httptest.NewRecorder()
	handler.ListLoans(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if gotFilter.ClientCPF == nil || *gotFilter.ClientCPF != "12345678900" {
		t.Fatalf("expected ClientCPF '12345678900' (digits only), got %+v", gotFilter.ClientCPF)
	}
	if gotFilter.PartnerCPF == nil || *gotFilter.PartnerCPF != "98765432100" {
		t.Fatalf("expected PartnerCPF '98765432100', got %+v", gotFilter.PartnerCPF)
	}

	// punctuation-only input has no digits -> nil filter
	req2 := httptest.NewRequest(http.MethodGet, "/admin/loans?clientCpf=...---", nil)
	rec2 := httptest.NewRecorder()
	handler.ListLoans(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 for punctuation-only cpf, got %d, body: %s", rec2.Code, rec2.Body.String())
	}
	if gotFilter.ClientCPF != nil {
		t.Fatalf("expected nil ClientCPF for punctuation-only input, got %+v", gotFilter.ClientCPF)
	}
}
