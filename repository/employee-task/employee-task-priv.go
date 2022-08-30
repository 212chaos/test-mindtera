package repo_employeetask

import (
	"fmt"
	"time"
)

func (e *EmployeeTaskRepositoryImpl) getColumnsNeededJoins() string {
	return `employee_tasks.id, employee_tasks.user_id, 
	employee_tasks.email, employee_tasks.due_date, 
	employee_tasks.assign_date, employee_tasks.task_type, 
	employee_tasks.task_related_id, employee_tasks.task_status, 
	employee_tasks.subscription_id, corporate_employee_entities.employee_name as "employee_name",
	corporate_employee_entities.department as "department"`
}

func (e *EmployeeTaskRepositoryImpl) getColumnsNeeded() string {
	return `id, user_id, email, due_date, assign_date, task_type, task_status, subscription_id`
}

func (e *EmployeeTaskRepositoryImpl) insertNewWellBeing(start, dueDate time.Time) string {
	oneDayBeforeMonth := dueDate.AddDate(0, -1, -7)
	_, end := e.commonSvc.GenerateFirstEndMonth(oneDayBeforeMonth)

	return fmt.Sprintf(
		`INSERT INTO employee_tasks( 
			email, 
			due_date, 
			task_type,
			user_id,
			task_status,
			created_at,
			record_flag,
			subscription_id,
			  created_by
		)
		SELECT  e.email,
				'%v' as due_date,
				e.task_type,
				e.user_id,
				'ACTIVE' as task_status,
				NOW() as created_at,
				'ACTIVE' as record_flag,
				e.subscription_id,
				'SCHEDULER' as created_by
		FROM employee_tasks e 
		LEFT JOIN corporate_subscription_entities c
		ON e.subscription_id = c.id
		WHERE e.task_type='WELL_BEING'
			AND c.end_period > NOW()
			AND (e.due_date >= '%v' AND e.due_date < '%v');`,
		dueDate.Format(time.RFC3339),
		end.Format(time.RFC3339),
		start.Format(time.RFC3339))
}

func (e *EmployeeTaskRepositoryImpl) updateWellBeingPreviousMonth(start, dueDate time.Time) string {
	oneDayBeforeMonth := dueDate.AddDate(0, -1, -7)
	_, end := e.commonSvc.GenerateFirstEndMonth(oneDayBeforeMonth)

	return fmt.Sprintf(`UPDATE employee_tasks 
	SET task_status='EXPIRED',
		record_flag='EXPIRED', 
		deleted_at = NOW() 
	WHERE task_type = 'WELL_BEING'
		AND task_status = 'ACTIVE'
		AND (due_date >= '%v' AND due_date < '%v');`,
		end.Format(time.RFC3339),
		start.Format(time.RFC3339))
}
