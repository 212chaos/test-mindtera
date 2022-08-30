package common

import (
	"os"
	"strconv"
	"strings"
)

type CONTEXT_KEY string

const (
	PUBLIC_KEY CONTEXT_KEY = "PUBLIC_KEY"
)

var (
	SERVICE_NAME                                     = os.Getenv("SERVICE_NAME")
	SERVICE_PORT                                     = os.Getenv("SERVICE_PORT")
	GRPC_PORT                                        = os.Getenv("GRPC_PORT")
	GIN_MODE                                         = os.Getenv("GIN_MODE")
	ENV                                              = os.Getenv("ENV")
	API_URL                                          = os.Getenv("API_URL")
	API_PUBLIC_URL                                   = os.Getenv("API_PUBLIC_URL")
	AUTH_METHOD                                      = os.Getenv("AUTH_METHOD")
	REDIS_AUTH_PREFIX                                = os.Getenv("REDIS_AUTH_PREFIX")
	ICON_BUCKET                                      = os.Getenv("ICON_BUCKET")
	ICON_OUTLINE_VALUE_FOLDER                        = os.Getenv("ICON_OUTLINE_VALUE_FOLDER")
	ICON_VALUE_FOLDER                                = os.Getenv("ICON_VALUE_FOLDER")
	CONSUMER_PUBLIC_PATH_PREFIX                      = os.Getenv("CONSUMER_PUBLIC_PATH_PREFIX")
	CONSUMER_SERVICE_URL                             = os.Getenv("CONSUMER_SERVICE_URL")
	CORPORATE_SERVICE_CALLBACK_URL                   = os.Getenv("CORPORATE_SERVICE_CALLBACK_URL")
	STATIC_AUTH                                      = os.Getenv("STATIC_AUTH")
	MAILJET_MAIL_TO_CORPORATE_SERVICE_BOOKING        = os.Getenv("MAILJET_MAIL_TO_CORPORATE_SERVICE_BOOKING")
	REGISTRATION_EMAIL_SUBJECT                       = strings.ReplaceAll(os.Getenv("REGISTRATION_EMAIL_SUBJECT"), `"`, "")
	ASSESSMENT_REMINDER_SUBJECT                      = strings.ReplaceAll(os.Getenv("ASSESSMENT_REMINDER_SUBJECT"), `"`, "")
	COMMON_TTL, _                                    = strconv.Atoi(os.Getenv("COMMON_TTL"))
	REGISTRATION_EMAIL_TEMPLATE, _                   = strconv.Atoi(os.Getenv("REGISTRATION_EMAIL_TEMPLATE"))
	ASSESSMENT_REMINDER_TEMPLATE_ID, _               = strconv.Atoi(os.Getenv("ASSESSMENT_REMINDER_TEMPLATE_ID"))
	MAILJET_CORPORATE_SERVICE_BOOKING_TEMPLATE_ID, _ = strconv.Atoi(os.Getenv("MAILJET_CORPORATE_SERVICE_BOOKING_TEMPLATE_ID"))
	HARD_DELETE_FLAG                                 = "HARD_DELETED"
	DELETE_FLAG                                      = "DELETED"
	DEACTIVATE_FLAG                                  = "DEACTIVATED"

	ACTIVE_RECORD_FLAG_QUERY                 = "(record_flag = 'ACTIVE')"
	ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE       = "(record_flag = 'ACTIVE') and hide = false"
	ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED = "(record_flag IN('ACTIVE','DEACTIVATED','UNREGISTERED'))"

	EMPLOYEE_TASK_BUCKET             = os.Getenv("EMPLOYEE_TASK_BUCKET")
	CORPORATE_SERVICE_BOOKING_BUCKET = os.Getenv("CORPORATE_SERVICE_BOOKING_BUCKET")

	DELETION_VALUE_WHITELIST = strings.Split(os.Getenv("DELETION_VALUE_WHITELIST"), ",")
)
