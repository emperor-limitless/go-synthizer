package synthizer

/*
#include <synthizer.h>
#include <synthizer_constants.h>
#include <stdlib.h>
*/
import "C"

type LogLevel C.int

type LoggingBackend C.int

const (
	LOG_LEVEL_ERROR LogLevel = C.SYZ_LOG_LEVEL_ERROR
	LOG_LEVEL_WARN = C.SYZ_LOG_LEVEL_WARN
	LOG_LEVEL_INFO = C.SYZ_LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG = C.SYZ_LOG_LEVEL_DEBUG
)

const (
	LOGGING_BACKEND_NONE LoggingBackend = C.SYZ_LOGGING_BACKEND_NONE
	LOGGING_BACKEND_STDERR = C.SYZ_LOGGING_BACKEND_STDERR
)