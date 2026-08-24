package input

// Pagination holds page/limit for list endpoints. Page is 1-based.
type Pagination struct {
	Page  int
	Limit int
}

// ListLoanFilter holds optional filters for listing loans.
// A nil pointer means "no filter" for that field.
type ListLoanFilter struct {
	PaymentStatus *string
	ClientName    *string
	PartnerName   *string
	ClientCPF     *string
	PartnerCPF    *string
}
