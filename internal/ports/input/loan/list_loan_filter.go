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
	Type          *int
	ClientName    *string
	PartnerName   *string

	AmountMin        *float64
	AmountMax        *float64
	AskValueMin      *float64
	AskValueMax      *float64
	ClientAmountMin  *float64
	ClientAmountMax  *float64
	GrossProfitMin   *float64
	GrossProfitMax   *float64
	ProfitMin        *float64
	ProfitMax        *float64
	PartnerAmountMin *float64
	PartnerAmountMax *float64
	OperationPctMin  *float64
	OperationPctMax  *float64
	PartnerPctMin    *float64
	PartnerPctMax    *float64
	NumberCardsMin   *int
	NumberCardsMax   *int
}
