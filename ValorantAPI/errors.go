package valorantapi

const (
	/* [ Errors ] */

	VALORANT_NO_VALID_RESPONSE = " No valid ValorantAPI response"
	VALORANT_API_UNAUTHORIZED  = " API Request was unauthorized"
	VALORANT_API_TIMEOUT       = " Too many ValorantAPI Requests"

	/* [ REGEX Errors ] */

	VALORANT_NO_REGION_SHARD = " Could not find region/shard data in log file"

	/* [ Authentication Errors ] */

	VALORANT_NO_LOCAL_LOCKFILE_FOUND             = " No Local Lockfile was found, RiotClient may be closed"
	VALORANT_AUTHENTICATION_INVALID_AUTH_COOKIES = " AuthCookies provided were invalid/nil"
	VALORANT_AUTHENTICATION_FAIL                 = " Authentication has failed"
)
