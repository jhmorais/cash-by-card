package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
)

const (
	maxLoansLimit = 100
)

func parseStringFilterParam(query url.Values, name string) *string {
	value := query.Get(name)
	if value == "" {
		return nil
	}
	return &value
}

// parseCpfFilterParam returns the param value reduced to digits only.
// Returns nil when absent, empty, or containing no digits.
func parseCpfFilterParam(query url.Values, name string) *string {
	value := query.Get(name)
	if value == "" {
		return nil
	}
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, value)
	if digits == "" {
		return nil
	}
	return &digits
}

// parseListLoansParams reads and validates the query params of GET /admin/loans.
func parseListLoansParams(r *http.Request) (*input.ListLoanFilter, *input.Pagination, error) {
	query := r.URL.Query()

	pagination := &input.Pagination{Page: 1, Limit: 10}

	if raw := query.Get("page"); raw != "" {
		page, err := strconv.Atoi(raw)
		if err != nil || page < 1 {
			return nil, nil, fmt.Errorf("invalid 'page' parameter, expected an integer >= 1")
		}
		pagination.Page = page
	}

	if raw := query.Get("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > maxLoansLimit {
			return nil, nil, fmt.Errorf("invalid 'limit' parameter, expected an integer between 1 and %d", maxLoansLimit)
		}
		pagination.Limit = limit
	}

	filter := &input.ListLoanFilter{}

	if status := query.Get("paymentStatus"); status != "" && status != "pending" && status != "paid" {
		return nil, nil, fmt.Errorf("invalid 'paymentStatus' parameter, expected 'pending' or 'paid'")
	}
	filter.PaymentStatus = parseStringFilterParam(query, "paymentStatus")

	filter.ClientName = parseStringFilterParam(query, "clientName")
	filter.PartnerName = parseStringFilterParam(query, "partnerName")
	filter.ClientCPF = parseCpfFilterParam(query, "clientCpf")
	filter.PartnerCPF = parseCpfFilterParam(query, "partnerCpf")

	return filter, pagination, nil
}
