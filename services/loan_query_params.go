package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

func parseIntFilterParam(query url.Values, name string) (*int, error) {
	value := query.Get(name)
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid '%s' parameter, expected an integer", name)
	}
	return &parsed, nil
}

func parseFloatFilterParam(query url.Values, name string) (*float64, error) {
	value := query.Get(name)
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid '%s' parameter, expected a number", name)
	}
	return &parsed, nil
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

	loanType, err := parseIntFilterParam(query, "type")
	if err != nil {
		return nil, nil, err
	}
	if loanType != nil && *loanType != 1 && *loanType != 2 {
		return nil, nil, fmt.Errorf("invalid 'type' parameter, expected 1 or 2")
	}
	filter.Type = loanType

	filter.ClientName = parseStringFilterParam(query, "clientName")
	filter.PartnerName = parseStringFilterParam(query, "partnerName")

	floatRanges := []struct {
		minName string
		maxName string
		minDst  **float64
		maxDst  **float64
	}{
		{"amountMin", "amountMax", &filter.AmountMin, &filter.AmountMax},
		{"askValueMin", "askValueMax", &filter.AskValueMin, &filter.AskValueMax},
		{"clientAmountMin", "clientAmountMax", &filter.ClientAmountMin, &filter.ClientAmountMax},
		{"grossProfitMin", "grossProfitMax", &filter.GrossProfitMin, &filter.GrossProfitMax},
		{"profitMin", "profitMax", &filter.ProfitMin, &filter.ProfitMax},
		{"partnerAmountMin", "partnerAmountMax", &filter.PartnerAmountMin, &filter.PartnerAmountMax},
		{"operationPercentMin", "operationPercentMax", &filter.OperationPctMin, &filter.OperationPctMax},
		{"partnerPercentMin", "partnerPercentMax", &filter.PartnerPctMin, &filter.PartnerPctMax},
	}
	for _, fr := range floatRanges {
		minValue, err := parseFloatFilterParam(query, fr.minName)
		if err != nil {
			return nil, nil, err
		}
		maxValue, err := parseFloatFilterParam(query, fr.maxName)
		if err != nil {
			return nil, nil, err
		}
		*fr.minDst = minValue
		*fr.maxDst = maxValue
	}

	filter.NumberCardsMin, err = parseIntFilterParam(query, "numberCardsMin")
	if err != nil {
		return nil, nil, err
	}
	filter.NumberCardsMax, err = parseIntFilterParam(query, "numberCardsMax")
	if err != nil {
		return nil, nil, err
	}

	return filter, pagination, nil
}
