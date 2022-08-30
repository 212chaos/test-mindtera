package repository_corporateprogram

// first param is user_id second param subscription id
// TODO NOTE: THIS IS CONTAIN TO FETCH DATA FROM PROGRAM TAGGING RELATION
func (c *CorporateProgramImpl) getSkippedProgramAmount() string {
	return `SELECT task_related_id as "program_id", count(1) as amount
			FROM employee_tasks
				WHERE 
					(task_status IN ('SKIPPED') 
					OR record_flag IN ('SKIPPED'))
					AND user_id = ?
					AND subscription_id = ?
					AND task_type = 'PROGRAM_RECOMMENDATION'
			GROUP BY task_related_id`
}

// first param is corporate id, user id, and subscription id
// TODO NOTE: THIS IS CONTAIN TO FETCH DATA FROM ENROLLED PROGRAM AND PROGRAM TAGGING RELATION
func (c *CorporateProgramImpl) getAvailableProgramAmount() string {
	return `SELECT program_id as task_related_id, count(1) as amountCount
			FROM program_tagging_relations 
				WHERE tagging_id IN 
				(
					Select program_tagging_id::text
					FROM corporate_value_entities
					WHERE id IN 
					(
						SELECT corporate_value_id
						FROM corporate_value_relations
						WHERE corporate_id = ?
					)
				)
				AND program_id NOT IN 
				(
					SELECT task_related_id::text 
					FROM employee_tasks
					WHERE 
						(task_status = 'FINISHED' 
						OR record_flag = 'FINISHED')
						AND user_id = ?
						AND subscription_id = ?
						AND task_type = 'PROGRAM_RECOMMENDATION'
				)
				AND program_id NOT IN
				(
					SELECT program_id::text
					FROM program_enrollment_history
					WHERE 
						status IN ('FINISH', 'FINISHED') 
						AND record_flag != 'DELETED' AND user_id = ?
				)
			GROUP BY program_id
			ORDER BY amountCount DESC;`
}
