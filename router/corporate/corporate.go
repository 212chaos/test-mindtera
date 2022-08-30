package corporateroute

type CorporateRoute interface {
	getCorporateStatus()
	upsertCorporateSubscription()

	Routes()
}
