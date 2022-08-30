package corporateemployeeroute

type CorporateEmployeeRoute interface {
	getCorporateEmployee()
	upsertCorporateEmployee()
	deleteCorporateEmployee()

	Routes()
}
