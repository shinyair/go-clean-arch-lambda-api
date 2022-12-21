package authentication

const (
	JwtHeader       string = "Authorization"
	JwtHeaderPrefix string = "Bearer "
)

// UserContext struct saves user information.
type UserContext struct {
	UserID        string `json:"userId"`
	UserName      string `json:"userName"`
	Locale        string `json:"locale"`
	ZoneID        string `json:"zoneId"`
	PermissionBit uint64 `json:"pbit"`
}

// DeepCopy
//
//	@receiver c
//	@return UserContext
func (c *UserContext) DeepCopy() UserContext {
	return UserContext{
		UserID:        c.UserID,
		UserName:      c.UserName,
		Locale:        c.Locale,
		ZoneID:        c.ZoneID,
		PermissionBit: c.PermissionBit,
	}
}

// AuthJwtClaim struct saves jwt claims.
type AuthJwtClaim struct {
	User      *UserContext `json:"userContext"`
	Issuer    string       `json:"iss"`
	Audience  string       `json:"aud"`
	Subject   string       `json:"sub"`
	IssuesAt  int64        `json:"issAt"`
	ExpiresAt int64        `json:"exp"`
}

// Valid necessary function to implement jwt.Claims
//
//	@receiver c
//	@return error always return <value>nil</value>. check detail after claim is returned by jwt package
func (c *AuthJwtClaim) Valid() error {
	return nil
}
