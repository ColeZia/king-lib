package framework

const (
	ENV_LOCAL = "local"
	ENV_DEV   = "dev"
	ENV_TEST  = "test"
	ENV_PRE   = "pre"
	ENV_PROD  = "prod"

	METADATA_GLOBAL_PREFIX     = "X-Md-Global-"
	METADATA_LOCAL_PREFIX      = "X-Md-Local-"
	METADATA_KEY_AUTH_TOKEN    = METADATA_GLOBAL_PREFIX + "Auth-Token"
	METADATA_KEY_CALL_SCENARIO = METADATA_GLOBAL_PREFIX + "Call-Scenario"

	METADATA_PRISM_PREFIX           = "X-Md-Analytics-"
	METADATA_KEY_PRISM_ACCESS_TOKEN = METADATA_PRISM_PREFIX + "Access-Token"

	METADATA_OPEN_PREFIX           = "X-Md-Open-"
	METADATA_KEY_OPEN_ACCESS_TOKEN = METADATA_OPEN_PREFIX + "Access-Token"

	BOSS_USER_TOKEN_SUFFIX = "Auth-Token-Op-User"
	BOSS_SVC_TOKEN_SUFFIX  = "Auth-Token-Svc"

	METADATA_KEY_GLOBAL_OP_USER_TOKEN = METADATA_GLOBAL_PREFIX + BOSS_USER_TOKEN_SUFFIX
	METADATA_KEY_GLOBAL_SVC_TOKEN     = METADATA_GLOBAL_PREFIX + BOSS_SVC_TOKEN_SUFFIX

	METADATA_KEY_GLOBAL_OP_USER_USERNAME      = METADATA_GLOBAL_PREFIX + "Op-User-Username"
	METADATA_KEY_GLOBAL_OP_USER_ID            = METADATA_GLOBAL_PREFIX + "Op-User-Id"
	METADATA_KEY_GLOBAL_OP_USER_NAME          = METADATA_GLOBAL_PREFIX + "Op-User-Name"
	METADATA_KEY_GLOBAL_OP_USER_OCEAN_TOKEN   = METADATA_GLOBAL_PREFIX + "Op-User-Ocean-Token"
	METADATA_KEY_GLOBAL_OP_USER_LOGIN_METHOD  = METADATA_GLOBAL_PREFIX + "Op-User-Login-Method"
	METADATA_KEY_GLOBAL_OP_USER_FEISHU_TOKEN  = METADATA_GLOBAL_PREFIX + "Op-User-Feishu-Token"
	METADATA_KEY_GLOBAL_OP_USER_CONSOLE_TOKEN = METADATA_GLOBAL_PREFIX + "Op-User-Console-Token"

	METADATA_KEY_LOCAL_OP_USER_TOKEN = METADATA_LOCAL_PREFIX + BOSS_USER_TOKEN_SUFFIX
	METADATA_KEY_LOCAL_SVC_TOKEN     = METADATA_LOCAL_PREFIX + BOSS_SVC_TOKEN_SUFFIX

	MDV_SERVICE_CALL_SCENARIO_UNKOWN  = "unkown"
	MDV_SERVICE_CALL_SCENARIO_GETEWAY = "gateway"
	MDV_SERVICE_CALL_SCENARIO_INNER   = "inner"

	AUTH_METHOD_SESSION = "SESSION"
	AUTH_METHOD_JWT     = "JWT"
	AUTH_METHOD_TOKEN   = "TOKEN"

	AUTH_PLATFORM_BOSS    = 1
	AUTH_PLATFORM_CONSOLE = 2

	FRAMEWORK_NAME_GO_KRATOS = "go-kratos"
	FRAMEWORK_NAME_GO_MICRO  = "go-micro"
	FRAMEWORK_NAME_GO_KIT    = "go-kit"
	FRAMEWORK_NAME_GIN       = "gin"

	K8sSMEnvPreName   = "pre"
	K8sSMEnvProdName  = "prod"
	K8sSMEnvAlphaName = "alpha"
)