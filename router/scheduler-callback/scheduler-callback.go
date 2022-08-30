package schedulercallbackroute

type SchedulerCallbackRoute interface {
	postAssessmentCallback()
	postExpiredSubscriptionCallback()

	Routes()
}
