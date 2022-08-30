package corporatevaluerelationroute

type CorporateValueRelationRoute interface {
	getCorporateValueRelation()
	upsertCorporateValueRelation()
	deleteCorporateValueRelation()

	Routes()
}
