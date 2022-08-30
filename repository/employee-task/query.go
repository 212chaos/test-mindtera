package repo_employeetask

import (
	"fmt"
	"time"

	"github.com/mindtera/corporate-service/entity"
)

// insertTaskEmployee
// createdBy refers to corporateID
func insertTaskEmployee(tasks []entity.EmployeeTaskEntity) string {
	var finalQuery string = `INSERT INTO public.employee_tasks ( user_id, email, assign_date, task_type, task_related_id, task_status, created_by) VALUES`
	var dynamicQuery string = `(`
	mappingLen := len(tasks)

	for i, v := range tasks {
		if i == mappingLen-1 {
			dynamicQuery = fmt.Sprintf("%s'%s', '%s', CAST('%s' AS DATE), '%s', '%s', '%s', '%s'", dynamicQuery, v.UserID, v.Email, v.AssignDate.Format(time.RFC3339), v.TaskType, v.TaskRelatedID, v.TaskStatus, v.CreatedBy)
		} else {
			dynamicQuery = fmt.Sprintf("%s'%s', '%s', CAST('%s' AS DATE), '%s', '%s', '%s', '%s'),(", dynamicQuery, v.UserID, v.Email, v.AssignDate.Format(time.RFC3339), v.TaskType, v.TaskRelatedID, v.TaskStatus, v.CreatedBy)
		}
	}
	dynamicQuery = fmt.Sprintf("%s)", dynamicQuery)
	return fmt.Sprintf("%s%s", finalQuery, dynamicQuery)
}

// insertProgramAssignmentHistory
func insertProgramAssignmentHistory(programIDArr []string, corporateID string, assignmentType entity.TypeAssigningMode, relatedID, createdBy string) string {
	var finalQuery string = `INSERT INTO public.employee_program_assignment_history ( program_id, corporate_id, assigning_type, related_id, created_by) VALUES`
	var dynamicQuery string = `(`
	mappingLen := len(programIDArr)

	for i, v := range programIDArr {
		if i == mappingLen-1 {
			dynamicQuery = fmt.Sprintf("%s'%s', '%s', '%s', '%s', '%s'", dynamicQuery, v, corporateID, assignmentType, relatedID, createdBy)
		} else {
			dynamicQuery = fmt.Sprintf("%s'%s', '%s', '%s', '%s', '%s'),(", dynamicQuery, v, corporateID, assignmentType, relatedID, createdBy)
		}
	}
	dynamicQuery = fmt.Sprintf("%s)", dynamicQuery)
	return fmt.Sprintf("%s%s", finalQuery, dynamicQuery)

}

// getAssignedProgram
// createdBy refers to corporateID
func getAssignedProgram(corporateID string, taskType entity.EMPLOYEE_TASK_TYPE, taskStatus entity.EMPLOYEE_TASK_STATUS, limit, offset int) string {
	return fmt.Sprintf(`SELECT
	employee_tasks.id AS "employee_tasks.id",
	employee_tasks.user_id AS "employee_tasks.user_id",
	employee_tasks.email AS "employee_tasks.email",
	employee_tasks.assign_date AS "employee_tasks.assign_date",
	employee_tasks.due_date AS "employee_tasks.due_date",
	employee_tasks.task_type AS "employee_tasks.task_type",
	employee_tasks.task_related_id AS "employee_tasks.task_related_id",
	employee_tasks.task_status AS "employee_tasks.task_status",
	employee_tasks.created_at AS "employee_tasks.created_at",
	employee_tasks.updated_at AS "employee_tasks.updated_at",
	employee_tasks.created_by AS "employee_tasks.created_by",
	employee_tasks.updated_by AS "employee_tasks.updated_by",
	employee_tasks.record_flag AS "employee_tasks.record_flag",
	employee_tasks.deleted_at AS "employee_tasks.deleted_at",
	employee_tasks.subscription_id AS "employee_tasks.subscription_id",
	corporate_employee_entities.id AS "corporate_employee_entities.id",
	corporate_employee_entities.public_user_id AS "corporate_employee_entities.public_user_id",
	corporate_employee_entities.corporate_id AS "corporate_employee_entities.corporate_id",
	corporate_employee_entities.department AS "corporate_employee_entities.department",
	corporate_employee_entities.role AS "corporate_employee_entities.role",
	corporate_employee_entities.email AS "corporate_employee_entities.email",
	corporate_employee_entities.record_flag AS "corporate_employee_entities.record_flag",
	corporate_employee_entities.created_at AS "corporate_employee_entities.created_at",
	corporate_employee_entities.updated_at AS "corporate_employee_entities.updated_at",
	corporate_employee_entities.deleted_at AS "corporate_employee_entities.deleted_at",
	corporate_employee_entities.created_by AS "corporate_employee_entities.created_by",
	corporate_employee_entities.updated_by AS "corporate_employee_entities.updated_by",
	corporate_employee_entities.employee_name AS "corporate_employee_entities.employee_name",
	corporate_employee_entities.corporate_subscription_id AS "corporate_employee_entities.corporate_subscription_id"
   FROM
	employee_tasks
   INNER JOIN corporate_employee_entities ON
	employee_tasks.user_id = corporate_employee_entities.public_user_id
   WHERE
	corporate_employee_entities.corporate_id = '%s'
	AND corporate_employee_entities.record_flag = 'ACTIVE'
	AND employee_tasks.task_type = '%s'
	AND employee_tasks.task_status  = '%s'
   ORDER BY
	employee_tasks.created_at DESC
   LIMIT %d OFFSET %d;`,
		corporateID, taskType, taskStatus, limit, offset,
	)
}

func countAssignedProgram(corporateID string, taskType entity.EMPLOYEE_TASK_TYPE, taskStatus entity.EMPLOYEE_TASK_STATUS) string {
	return fmt.Sprintf(`SELECT
		count(*)
	FROM
		employee_tasks
	INNER JOIN corporate_employee_entities ON
		employee_tasks.user_id = corporate_employee_entities.public_user_id
	WHERE
		corporate_employee_entities.corporate_id = '%s'
		AND corporate_employee_entities.record_flag = 'ACTIVE'
		AND employee_tasks.task_type = '%s'
		AND employee_tasks.task_status = '%s'`,
		corporateID, taskType, taskStatus,
	)
}

// searchPrograms
func searchPrograms(keyword string) string {
	return fmt.Sprintf(`SELECT m_programs.id, m_programs.name, m_programs.category, m_programs.sub_category, m_programs.description, m_programs.objective, m_programs.thumbnail_url, m_programs.rating, m_programs.duration, m_programs.subscription_type, m_programs.intro_video_id, m_programs.created_at, m_programs.updated_at, m_programs.created_by, m_programs.updated_by, m_programs.record_flag, m_programs.coach_id, m_programs.status, m_programs.content_type, m_programs.deleted_at, m_programs.moderator_id, m_programs.embedded_url, m_programs.published_at FROM m_programs WHERE upper(m_programs."name") LIKE '%%%s%%' AND record_flag = 'ACTIVE' AND status = 'PUBLISHED' ORDER BY m_programs.name LIMIT 5`,
		keyword)
}

// getPrograms
func getPrograms() string {
	return `SELECT mp.id AS "mp.id", mp."name" AS "mp.name", mp.category AS "mp.category", mp.sub_category AS "mp.sub_category", mp.description AS "mp.description", mp.objective AS "mp.objective", mp.thumbnail_url AS "mp.thumbnail_url", mp.rating AS "mp.rating", mp.duration AS "mp.duration", mp.subscription_type AS "mp.subscription_type", mp.intro_video_id AS "mp.intro_video_id", mp.created_at AS "mp.created_at", mp.updated_at AS "mp.updated_at", mp.created_by AS "mp.created_by", mp.updated_by AS "mp.updated_by", mp.record_flag AS "mp.record_flag", mp.coach_id AS "mp.coach_id", mp.status AS "mp.status", mp.content_type AS "mp.content_type", mc.id AS "mc.id", mc."name" AS "mc.name", mc.description AS "mc.description", mc.profile_picture AS "mc.profile_picture", mc.created_at AS "mc.created_at", mc.updated_at AS "mc.updated_at", mc.created_by AS "mc.created_by", mc.updated_by AS "mc.updated_by", mc.record_flag AS "mc.record_flag", mc."role" AS "mc.role", mc.origin AS "mc.origin", mc."level" AS "mc.level" FROM m_programs mp INNER JOIN m_coaches mc ON mp.coach_id = mc.id WHERE mp.status = 'PUBLISHED' AND mp.record_flag = 'ACTIVE' AND mc.record_flag = 'ACTIVE' ORDER BY mp.created_at DESC LIMIT 5`
}

func getCorporateServices(includeWebinar, includeExercise, includeCoaching bool, limit, offset int) string {
	initQuery := `SELECT m_corporate_services.id AS "m_corporate_services.id", m_corporate_services.name AS "m_corporate_services.name", m_corporate_services.description AS "m_corporate_services.description", m_corporate_services.long_description AS "m_corporate_services.long_description", m_corporate_services.landing_url AS "m_corporate_services.landing_url", m_corporate_services.image_url AS "m_corporate_services.image_url", m_corporate_services.service_type AS "m_corporate_services.service_type", m_corporate_services.created_at AS "m_corporate_services.created_at", m_corporate_services.updated_at AS "m_corporate_services.updated_at", m_corporate_services.created_by AS "m_corporate_services.created_by", m_corporate_services.updated_by AS "m_corporate_services.updated_by", m_corporate_services.record_flag AS "m_corporate_services.record_flag", m_corporate_services.deleted_at AS "m_corporate_services.deleted_at" FROM m_corporate_services WHERE m_corporate_services.record_flag = 'ACTIVE'`

	dynamicQuery := ""
	if includeWebinar {
		dynamicQuery = fmt.Sprintf("AND (m_corporate_services.service_type = '%s'", entity.WEBINAR_TYPE)
	}

	if includeExercise {
		if dynamicQuery == "" {
			dynamicQuery = fmt.Sprintf("AND (m_corporate_services.service_type = '%s'", entity.EXERCISE_TYPE)
		} else {
			dynamicQuery = fmt.Sprintf("%s OR m_corporate_services.service_type = '%s'", dynamicQuery, entity.EXERCISE_TYPE)
		}
	}

	if includeCoaching {
		if dynamicQuery == "" {
			dynamicQuery = fmt.Sprintf("AND (m_corporate_services.service_type = '%s'", entity.COACHING_TYPE)
		} else {
			dynamicQuery = fmt.Sprintf("%s OR m_corporate_services.service_type = '%s'", dynamicQuery, entity.COACHING_TYPE)
		}
	}

	if dynamicQuery != "" {
		dynamicQuery += ")"
	}

	orderLimitOffset := fmt.Sprintf("ORDER BY m_corporate_services.created_at DESC LIMIT %d OFFSET %d", limit, offset)

	var finalQuery string
	if dynamicQuery != "" {
		finalQuery = fmt.Sprintf("%s %s %s", initQuery, dynamicQuery, orderLimitOffset)
	} else {
		finalQuery = fmt.Sprintf("%s %s", initQuery, orderLimitOffset)
	}

	return finalQuery
}

func countGetCorporateServices(includeWebinar, includeExercise, includeCoaching bool) string {
	initQuery := `SELECT count(*) FROM m_corporate_services WHERE m_corporate_services.record_flag = 'ACTIVE'`

	dynamicQuery := ""
	if includeWebinar {
		dynamicQuery = fmt.Sprintf("AND (m_corporate_services.service_type = '%s'", entity.WEBINAR_TYPE)
	}

	if includeExercise {
		if dynamicQuery == "" {
			dynamicQuery = fmt.Sprintf("AND (m_corporate_services.service_type = '%s'", entity.EXERCISE_TYPE)
		} else {
			dynamicQuery = fmt.Sprintf("%s OR m_corporate_services.service_type = '%s'", dynamicQuery, entity.EXERCISE_TYPE)
		}
	}

	if includeCoaching {
		if dynamicQuery == "" {
			dynamicQuery = fmt.Sprintf("AND (m_corporate_services.service_type = '%s'", entity.COACHING_TYPE)
		} else {
			dynamicQuery = fmt.Sprintf("%s OR m_corporate_services.service_type = '%s'", dynamicQuery, entity.COACHING_TYPE)
		}
	}

	if dynamicQuery != "" {
		dynamicQuery += ")"
	}

	var finalQuery string
	if dynamicQuery != "" {
		finalQuery = fmt.Sprintf("%s %s", initQuery, dynamicQuery)
	} else {
		finalQuery = initQuery
	}

	return finalQuery
}

func getCorporateService(id string) string {
	return fmt.Sprintf(`SELECT
		m_corporate_services.id AS "m_corporate_services.id",
		m_corporate_services.name AS "m_corporate_services.name",
		m_corporate_services.description AS "m_corporate_services.description",
		m_corporate_services.long_description AS "m_corporate_services.long_description",
		m_corporate_services.landing_url AS "m_corporate_services.landing_url",
		m_corporate_services.image_url AS "m_corporate_services.image_url",
		m_corporate_services.service_type AS "m_corporate_services.service_type",
		m_corporate_services.created_at AS "m_corporate_services.created_at",
		m_corporate_services.updated_at AS "m_corporate_services.updated_at",
		m_corporate_services.created_by AS "m_corporate_services.created_by",
		m_corporate_services.updated_by AS "m_corporate_services.updated_by",
		m_corporate_services.record_flag AS "m_corporate_services.record_flag",
		m_corporate_services.deleted_at AS "m_corporate_services.deleted_at"
	FROM
		m_corporate_services
	WHERE
		id = '%s'
		AND record_flag = 'ACTIVE'`,
		id,
	)
}

// insertTaskEmployee
// createdBy refers to corporateID
func insertCorporateServiceBooking(bookings []entity.CorporateServiceBooking) string {
	var finalQuery string = `INSERT INTO corporate_service_bookings ( service_id, corporate_id, scheduled_at, status, created_by) VALUES`
	var dynamicQuery string = `(`
	mappingLen := len(bookings)

	for i, v := range bookings {
		if i == mappingLen-1 {
			dynamicQuery = fmt.Sprintf("%s'%s', '%s', CAST('%s' AS TIMESTAMP), '%s', '%s'", dynamicQuery, v.ServiceID, v.CorporateID, v.ScheduledAt.Format(time.RFC3339), v.Status, v.CorporateID)
		} else {
			dynamicQuery = fmt.Sprintf("%s'%s', '%s', CAST('%s' AS TIMESTAMP), '%s', '%s'),(", dynamicQuery, v.ServiceID, v.CorporateID, v.ScheduledAt.Format(time.RFC3339), v.Status, v.CorporateID)
		}
	}
	dynamicQuery = fmt.Sprintf("%s)", dynamicQuery)
	return fmt.Sprintf("%s%s", finalQuery, dynamicQuery)
}

func getCorporateServiceBooking(corporateID string, limit, offset int) string {
	return fmt.Sprintf(`SELECT corporate_service_bookings.id AS "corporate_service_bookings.id", corporate_service_bookings.service_id AS "corporate_service_bookings.service_id", corporate_service_bookings.corporate_id AS "corporate_service_bookings.corporate_id", corporate_service_bookings.scheduled_at AS "corporate_service_bookings.scheduled_at", corporate_service_bookings.status AS "corporate_service_bookings.status", corporate_service_bookings.created_at AS "corporate_service_bookings.created_at", corporate_service_bookings.updated_at AS "corporate_service_bookings.updated_at", corporate_service_bookings.created_by AS "corporate_service_bookings.created_by", corporate_service_bookings.updated_by AS "corporate_service_bookings.updated_by", corporate_service_bookings.record_flag AS "corporate_service_bookings.record_flag", corporate_service_bookings.deleted_at AS "corporate_service_bookings.deleted_at", m_corporate_services.id AS "m_corporate_services.id", m_corporate_services.name AS "m_corporate_services.name", m_corporate_services.description AS "m_corporate_services.description", m_corporate_services.landing_url AS "m_corporate_services.landing_url", m_corporate_services.image_url AS "m_corporate_services.image_url", m_corporate_services.service_type AS "m_corporate_services.service_type", m_corporate_services.created_at AS "m_corporate_services.created_at", m_corporate_services.updated_at AS "m_corporate_services.updated_at", m_corporate_services.created_by AS "m_corporate_services.created_by", m_corporate_services.updated_by AS "m_corporate_services.updated_by", m_corporate_services.record_flag AS "m_corporate_services.record_flag", m_corporate_services.deleted_at AS "m_corporate_services.deleted_at", m_corporate_services.long_description AS "m_corporate_services.long_description" FROM corporate_service_bookings INNER JOIN m_corporate_services ON m_corporate_services.id = corporate_service_bookings.service_id WHERE corporate_service_bookings.corporate_id = '%s' AND corporate_service_bookings.record_flag = 'ACTIVE' ORDER BY corporate_service_bookings.created_at DESC LIMIT %d OFFSET %d`,
		corporateID, limit, offset,
	)
}

func countCorporateServiceBooking(corporateID string) string {
	return fmt.Sprintf(`SELECT count(*) FROM corporate_service_bookings WHERE corporate_service_bookings.corporate_id ='%s' AND corporate_service_bookings.record_flag = 'ACTIVE'`,
		corporateID,
	)
}

// searchCorporateService
func searchCorporateService(keyword string) string {
	return fmt.Sprintf(`SELECT
		m_corporate_services.id AS "m_corporate_services.id",
		m_corporate_services.name AS "m_corporate_services.name",
		m_corporate_services.description AS "m_corporate_services.description",
		m_corporate_services.long_description AS "m_corporate_services.long_description",
		m_corporate_services.landing_url AS "m_corporate_services.landing_url",
		m_corporate_services.image_url AS "m_corporate_services.image_url",
		m_corporate_services.service_type AS "m_corporate_services.service_type",
		m_corporate_services.created_at AS "m_corporate_services.created_at",
		m_corporate_services.updated_at AS "m_corporate_services.updated_at",
		m_corporate_services.created_by AS "m_corporate_services.created_by",
		m_corporate_services.updated_by AS "m_corporate_services.updated_by",
		m_corporate_services.record_flag AS "m_corporate_services.record_flag",
		m_corporate_services.deleted_at AS "m_corporate_services.deleted_at"
	FROM
		m_corporate_services
	WHERE
		upper(m_corporate_services."name") LIKE '%%%s%%'
		AND record_flag = 'ACTIVE'
	ORDER BY
		m_corporate_services.name
	LIMIT 5`,
		keyword)
}

// GetEmployeeByUserID
func getEmployeeByUserID(userID string) string {
	return fmt.Sprintf(`SELECT corporate_department_entities.id AS "corporate_department_entities.id", corporate_department_entities.code AS "corporate_department_entities.code", corporate_department_entities.name AS "corporate_department_entities.name", corporate_department_entities.corporate_id AS "corporate_department_entities.corporate_id", corporate_department_entities.created_at AS "corporate_department_entities.created_at", corporate_department_entities.updated_at AS "corporate_department_entities.updated_at", corporate_department_entities.deleted_at AS "corporate_department_entities.deleted_at", corporate_department_entities.created_by AS "corporate_department_entities.created_by", corporate_department_entities.updated_by AS "corporate_department_entities.updated_by", corporate_department_entities.record_flag AS "corporate_department_entities.record_flag", corporate_employee_entities.id AS "corporate_employee_entities.id", corporate_employee_entities.public_user_id AS "corporate_employee_entities.public_user_id", corporate_employee_entities.corporate_id AS "corporate_employee_entities.corporate_id", corporate_employee_entities.department AS "corporate_employee_entities.department", corporate_employee_entities.role AS "corporate_employee_entities.role", corporate_employee_entities.email AS "corporate_employee_entities.email", corporate_employee_entities.record_flag AS "corporate_employee_entities.record_flag", corporate_employee_entities.created_at AS "corporate_employee_entities.created_at", corporate_employee_entities.updated_at AS "corporate_employee_entities.updated_at", corporate_employee_entities.deleted_at AS "corporate_employee_entities.deleted_at", corporate_employee_entities.created_by AS "corporate_employee_entities.created_by", corporate_employee_entities.updated_by AS "corporate_employee_entities.updated_by", corporate_employee_entities.employee_name AS "corporate_employee_entities.employee_name", corporate_employee_entities.corporate_subscription_id AS "corporate_employee_entities.corporate_subscription_id" FROM corporate_employee_entities INNER JOIN corporate_department_entities ON corporate_employee_entities.department = corporate_department_entities.code WHERE corporate_employee_entities.public_user_id = '%s' AND corporate_department_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.record_flag = 'ACTIVE' ORDER BY corporate_employee_entities.created_at DESC LIMIT 1`,
		userID)
}

// GetEmployeeByCorporateID
// Get employee by corporate_id
func getEmployeeByCorporateID(corporateID string) string {
	return fmt.Sprintf(`SELECT corporate_department_entities.id AS "corporate_department_entities.id", corporate_department_entities.code AS "corporate_department_entities.code", corporate_department_entities.name AS "corporate_department_entities.name", corporate_department_entities.corporate_id AS "corporate_department_entities.corporate_id", corporate_department_entities.created_at AS "corporate_department_entities.created_at", corporate_department_entities.updated_at AS "corporate_department_entities.updated_at", corporate_department_entities.deleted_at AS "corporate_department_entities.deleted_at", corporate_department_entities.created_by AS "corporate_department_entities.created_by", corporate_department_entities.updated_by AS "corporate_department_entities.updated_by", corporate_department_entities.record_flag AS "corporate_department_entities.record_flag", corporate_employee_entities.id AS "corporate_employee_entities.id", corporate_employee_entities.public_user_id AS "corporate_employee_entities.public_user_id", corporate_employee_entities.corporate_id AS "corporate_employee_entities.corporate_id", corporate_employee_entities.department AS "corporate_employee_entities.department", corporate_employee_entities.role AS "corporate_employee_entities.role", corporate_employee_entities.email AS "corporate_employee_entities.email", corporate_employee_entities.record_flag AS "corporate_employee_entities.record_flag", corporate_employee_entities.created_at AS "corporate_employee_entities.created_at", corporate_employee_entities.updated_at AS "corporate_employee_entities.updated_at", corporate_employee_entities.deleted_at AS "corporate_employee_entities.deleted_at", corporate_employee_entities.created_by AS "corporate_employee_entities.created_by", corporate_employee_entities.updated_by AS "corporate_employee_entities.updated_by", corporate_employee_entities.employee_name AS "corporate_employee_entities.employee_name", corporate_employee_entities.corporate_subscription_id AS "corporate_employee_entities.corporate_subscription_id" FROM corporate_employee_entities INNER JOIN corporate_department_entities ON corporate_employee_entities.department = corporate_department_entities.code WHERE corporate_department_entities.corporate_id = '%s' AND corporate_department_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.record_flag = 'ACTIVE' ORDER BY corporate_employee_entities.created_at DESC`,
		corporateID)
}

// GetEmployeeByCorporateIDDepartmentCode
// Get employee list where department, corporate_id
func getEmployeeByCorporateIDDepartmentCode(corporateID, DepartmentCode string) string {
	return fmt.Sprintf(`SELECT corporate_department_entities.id AS "corporate_department_entities.id", corporate_department_entities.code AS "corporate_department_entities.code", corporate_department_entities.name AS "corporate_department_entities.name", corporate_department_entities.corporate_id AS "corporate_department_entities.corporate_id", corporate_department_entities.created_at AS "corporate_department_entities.created_at", corporate_department_entities.updated_at AS "corporate_department_entities.updated_at", corporate_department_entities.deleted_at AS "corporate_department_entities.deleted_at", corporate_department_entities.created_by AS "corporate_department_entities.created_by", corporate_department_entities.updated_by AS "corporate_department_entities.updated_by", corporate_department_entities.record_flag AS "corporate_department_entities.record_flag", corporate_employee_entities.id AS "corporate_employee_entities.id", corporate_employee_entities.public_user_id AS "corporate_employee_entities.public_user_id", corporate_employee_entities.corporate_id AS "corporate_employee_entities.corporate_id", corporate_employee_entities.department AS "corporate_employee_entities.department", corporate_employee_entities.role AS "corporate_employee_entities.role", corporate_employee_entities.email AS "corporate_employee_entities.email", corporate_employee_entities.record_flag AS "corporate_employee_entities.record_flag", corporate_employee_entities.created_at AS "corporate_employee_entities.created_at", corporate_employee_entities.updated_at AS "corporate_employee_entities.updated_at", corporate_employee_entities.deleted_at AS "corporate_employee_entities.deleted_at", corporate_employee_entities.created_by AS "corporate_employee_entities.created_by", corporate_employee_entities.updated_by AS "corporate_employee_entities.updated_by", corporate_employee_entities.employee_name AS "corporate_employee_entities.employee_name", corporate_employee_entities.corporate_subscription_id AS "corporate_employee_entities.corporate_subscription_id" FROM corporate_employee_entities INNER JOIN corporate_department_entities ON corporate_employee_entities.department = corporate_department_entities.code WHERE corporate_department_entities.corporate_id = '%s' AND corporate_department_entities.code = '%s' AND corporate_department_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.record_flag = 'ACTIVE' ORDER BY corporate_employee_entities.created_at DESC`,
		corporateID, DepartmentCode)
}

// GetProgramAssignmentHistory
func getProgramAssignmentHistory(corporateID string, limit, offset int) string {
	return fmt.Sprintf(`SELECT employee_program_assignment_history.id AS "employee_program_assignment_history.id", employee_program_assignment_history.program_id AS "employee_program_assignment_history.program_id", employee_program_assignment_history.corporate_id AS "employee_program_assignment_history.corporate_id", employee_program_assignment_history.assigning_type AS "employee_program_assignment_history.assigning_type", employee_program_assignment_history.related_id AS "employee_program_assignment_history.related_id", employee_program_assignment_history.created_at AS "employee_program_assignment_history.created_at", employee_program_assignment_history.updated_at AS "employee_program_assignment_history.updated_at", employee_program_assignment_history.created_by AS "employee_program_assignment_history.created_by", employee_program_assignment_history.updated_by AS "employee_program_assignment_history.updated_by", employee_program_assignment_history.record_flag AS "employee_program_assignment_history.record_flag", employee_program_assignment_history.deleted_at AS "employee_program_assignment_history.deleted_at", m_programs.id AS "m_programs.id", m_programs.name AS "m_programs.name", m_programs.category AS "m_programs.category", m_programs.sub_category AS "m_programs.sub_category", m_programs.description AS "m_programs.description", m_programs.objective AS "m_programs.objective", m_programs.thumbnail_url AS "m_programs.thumbnail_url", m_programs.rating AS "m_programs.rating", m_programs.duration AS "m_programs.duration", m_programs.subscription_type AS "m_programs.subscription_type", m_programs.intro_video_id AS "m_programs.intro_video_id", m_programs.created_at AS "m_programs.created_at", m_programs.updated_at AS "m_programs.updated_at", m_programs.created_by AS "m_programs.created_by", m_programs.updated_by AS "m_programs.updated_by", m_programs.record_flag AS "m_programs.record_flag", m_programs.coach_id AS "m_programs.coach_id", m_programs.status AS "m_programs.status", m_programs.content_type AS "m_programs.content_type", m_programs.deleted_at AS "m_programs.deleted_at", m_programs.moderator_id AS "m_programs.moderator_id", m_programs.embedded_url AS "m_programs.embedded_url", m_programs.published_at AS "m_programs.published_at" FROM employee_program_assignment_history INNER JOIN m_programs ON employee_program_assignment_history.program_id = m_programs.id WHERE corporate_id = '%s' AND employee_program_assignment_history.record_flag = 'ACTIVE' AND m_programs.record_flag = 'ACTIVE' AND m_programs.status = 'PUBLISHED' ORDER BY employee_program_assignment_history.created_at DESC LIMIT %d OFFSET %d`,
		corporateID, limit, offset,
	)
}

func countAssignedProgramHistory(corporateID string) string {
	return fmt.Sprintf(`SELECT count(*) FROM employee_program_assignment_history epah WHERE epah.corporate_id = '%s' AND epah.record_flag = 'ACTIVE'`,
		corporateID)
}
