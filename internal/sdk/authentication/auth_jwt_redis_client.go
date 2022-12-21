package authentication

// TODO:
// type AuthJwtRedisClient struct {
// 	publicKeyParam  string
// 	privateKeyParam string
// 	// no additional cost
// 	ssmClient   *ssm.SSM
// 	redisClient *redis.Client
// 	// have additional cost, or can use aws kms to issue&verify
// 	// secretClient     *awssecret.SecretsManager
// }

// // NewAuthJwtRedisClient
// //
// //	@param publicKeyParam
// //	@param privateKeyParam
// //	@param ssmClient
// //	@param redisClient
// //	@return *AuthJwtRedisClient
// func NewAuthJwtRedisClient(
// 	publicKeyParam string,
// 	privateKeyParam string,
// 	ssmClient *ssm.SSM,
// 	redisClient *redis.Client,
// ) *AuthJwtRedisClient {
// 	return &AuthJwtRedisClient{
// 		publicKeyParam:  publicKeyParam,
// 		privateKeyParam: privateKeyParam,
// 		ssmClient:       ssmClient,
// 		redisClient:     redisClient,
// 	}
// }

// func (c *AuthJwtRedisClient) Issue(ctx context.Context, claim *AuthJwtClaim) (string, error) {
// 	if claim == nil {
// 		return "", ErrInvalidClaim
// 	}
// 	privateKey, err := c.getSSMParam(ctx, c.privateKeyParam)
// 	if err != nil {
// 		return "", errors.Wrap(err, "failed to get private key")
// 	}
// 	rsakey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
// 	if err != nil {
// 		rootErr := errors.New(err.Error())
// 		return "", errors.Wrap(rootErr, "failed to parse private key")
// 	}
// 	curr := time.Now()
// 	claim.IssuesAt = curr.UnixNano()
// 	claim.ExpiresAt = curr.Add(ExpireDuration).UnixNano()
// 	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
// 	jwt, err := t.SignedString(rsakey)
// 	if err != nil {
// 		rootErr := errors.New(err.Error())
// 		return "", errors.Wrap(rootErr, "failed to issue jwt")
// 	}
// 	return jwt, nil
// }

// func (c *AuthJwtRedisClient) Verify(ctx context.Context, tokenStr string) (*AuthJwtClaim, error) {
// 	if len(tokenStr) == 0 {
// 		return nil, ErrInvalidJwt
// 	}
// 	blocked, err := c.isBlocked(ctx, tokenStr)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to check block status")
// 	}
// 	if blocked {
// 		return nil, ErrBlockedClaim
// 	}
// 	claim, err := c.parseJwt(ctx, tokenStr)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to prase jwt as claim")
// 	}
// 	if claim == nil || claim.User == nil {
// 		return nil, ErrInvalidClaim
// 	}
// 	if claim.ExpiresAt <= time.Now().UnixNano() {
// 		return nil, ErrExpiredClaim
// 	}
// 	return claim, nil
// }

// func (c *AuthJwtRedisClient) Block(ctx context.Context, tokenStr string) error {
// 	if len(tokenStr) == 0 {
// 		return nil
// 	}
// 	if c.redisClient == nil {
// 		return ErrBadClient
// 	}
// 	claim, err := c.parseJwt(ctx, tokenStr)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to parse jwt as claim")
// 	}
// 	if claim == nil {
// 		return errors.Wrapf(ErrInvalidJwt, "token: %s", tokenStr)
// 	}
// 	curr := time.Now().UnixNano()
// 	if claim.ExpiresAt <= curr {
// 		// already expired, no need to block it
// 		return nil
// 	}
// 	cmd := c.redisClient.Set(ctx, tokenStr, claim, time.Duration(claim.ExpiresAt-curr))
// 	_, err = cmd.Result()
// 	if err != nil {
// 		rootErr := errors.New(err.Error())
// 		return errors.Wrapf(rootErr, "failed to add token in redis: %s", tokenStr)
// 	}
// 	return nil
// }

// // isBlocked
// //
// //	@receiver c
// //	@param ctx
// //	@param tokenStr
// //	@return bool
// //	@return error
// func (c *AuthJwtRedisClient) isBlocked(ctx context.Context, tokenStr string) (bool, error) {
// 	if len(tokenStr) == 0 {
// 		return true, nil
// 	}
// 	if c.redisClient == nil {
// 		return false, ErrBadClient
// 	}
// 	cmd := c.redisClient.Get(ctx, tokenStr)
// 	_, err := cmd.Result()
// 	if err == nil {
// 		return true, nil
// 	}
// 	if errors.Is(err, redis.Nil) {
// 		return false, nil
// 	}
// 	rootErr := errors.New(err.Error())
// 	return false, errors.Wrapf(rootErr, "failed to get token from redis: %s", tokenStr)
// }

// // parseJwt
// //
// //	@receiver c
// //	@param ctx
// //	@param tokenStr
// //	@return *AuthJwtClaim
// //	@return error
// func (c *AuthJwtRedisClient) parseJwt(ctx context.Context, tokenStr string) (*AuthJwtClaim, error) {
// 	publicKey, err := c.getSSMParam(ctx, c.publicKeyParam)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to get public key")
// 	}
// 	token, err := jwt.ParseWithClaims(
// 		tokenStr,
// 		&AuthJwtClaim{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			//nolint:wrapcheck
// 			return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
// 		},
// 	)
// 	if err != nil {
// 		rootErr := errors.New(err.Error())
// 		return nil, errors.Wrapf(
// 			rootErr,
// 			"failed to parse token: %s, by public secret key: %s",
// 			tokenStr, c.publicKeyParam)
// 	}
// 	claim, ok := token.Claims.(*AuthJwtClaim)
// 	if !ok {
// 		return nil, ErrInvalidClaim
// 	}
// 	return claim, nil
// }

// // getSSMParam
// //
// //	@receiver c
// //	@param ctx
// //	@param paramName
// //	@return string
// //	@return error
// func (c *AuthJwtRedisClient) getSSMParam(ctx context.Context, paramName string) (string, error) {
// 	if paramName == "" {
// 		return "", errors.Wrapf(ErrInvalidJwt, "invalid key param name: %s", paramName)
// 	}
// 	if c.ssmClient == nil {
// 		return "", errors.Wrap(ErrBadClient, "no ssm client found")
// 	}
// 	input := &ssm.GetParameterInput{
// 		Name:           aws.String(paramName),
// 		WithDecryption: aws.Bool(false),
// 	}
// 	output, err := c.ssmClient.GetParameterWithContext(ctx, input)
// 	if err != nil {
// 		return "", errors.Wrap(ErrInvalidJwt, err.Error())
// 	}
// 	if output == nil || output.Parameter == nil || output.Parameter.Value == nil {
// 		return "", errors.Wrap(ErrInvalidJwt, "no ssm param found")
// 	}
// 	value := *output.Parameter.Value
// 	return value, nil
// }

// // getSecretKey
// //
// //	@receiver c
// //	@param ctx
// //	@param secretKey
// //	@return string
// //	@return error
// // func (c *AuthJwtRedisClient) getSecretKey(ctx context.Context, secretKey string) (string, error) {
// // 	if secretKey == "" || c.secretName == "" {
// // 		return "", errors.Wrapf(
// // 			ErrInvalidSecretManagerKey,
// // 			"secret key: %s, secret name: %s",
// // 			secretKey, c.secretName)
// // 	}
// // 	if c.secretClient == nil {
// // 		return "", errors.Wrap(ErrBadClient, "no secret client")
// // 	}
// // 	result, err := c.secretClient.GetSecretValueWithContext(ctx, &awssecret.GetSecretValueInput{
// // 		SecretId: aws.String(c.secretName),
// // 	})
// // 	if err != nil {
// // 		rootErr := errors.New(err.Error())
// // 		return "", errors.Wrapf(rootErr, "failed to get secret value by secret name: %s", c.secretName)
// // 	}
// // 	jsonMap := map[string]string{}
// // 	err = json.Unmarshal([]byte(*result.SecretString), &jsonMap)
// // 	if err != nil {
// // 		rootErr := errors.New(err.Error())
// // 		return "", errors.Wrapf(rootErr, "failed to unmarshal secret string as map: %s", *result.SecretString)
// // 	}
// // 	secret, ok := jsonMap[secretKey]
// // 	if !ok {
// // 		return "", errors.Wrapf(
// // 			ErrInvalidSecretManagerKey,
// // 			"no value found by secret key: %s, secret name: %s",
// // 			secretKey, c.secretName)
// // 	}
// // 	decoded, err := base64.StdEncoding.DecodeString(secret)
// // 	if err != nil {
// // 		rootErr := errors.New(err.Error())
// // 		return "", errors.Wrapf(rootErr, "cannot decode secret value: %s", secret)
// // 	}
// // 	return string(decoded), nil
// // }
