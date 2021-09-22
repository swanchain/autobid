package constants

const (
	SWAN_PAYMENT_ABI_JSON = "on-chain/contracts/abi/SwanPayment.json"
	DEFAULT_SELECT_LIMIT  = "100"

	URL_EVENT_PREFIX = "events"

	HTTP_STATUS_SUCCESS = "success"
	HTTP_STATUS_FAIL    = "fail"
	HTTP_STATUS_ERROR   = "error"

	HTTP_CODE_200_OK                    = "200" //http.StatusOk
	HTTP_CODE_400_BAD_REQUEST           = "400" //http.StatusBadRequest
	HTTP_CODE_401_UNAUTHORIZED          = "401" //http.StatusUnauthorized
	HTTP_CODE_500_INTERNAL_SERVER_ERROR = "500" //http.StatusInternalServerError

	URL_HOST_GET_COMMON       = "/common"
	URL_HOST_GET_HOST_INFO    = "/host/info"
	URL_HOST_GET_HEALTH_CHECK = "/health/check"

	TASK_STATUS_CREATED         = "Created"
	TASK_STATUS_ASSIGNED        = "Assigned"
	TASK_STATUS_ACTION_REQUIRED = "ActionRequired"
	TASK_STATUS_CANCELLED       = "Cancelled"
	TASK_STATUS_CLOSED          = "Closed"
	TASK_STATUS_COMPLETED       = "Completed"
	TASK_STATUS_EXPIRED         = "Expired"

	TASK_TYPE_REGULAR  = "regular"
	TASK_TYPE_VERIFIED = "verified"

	MINER_STATUS_ACTIVE = "Active"
)
