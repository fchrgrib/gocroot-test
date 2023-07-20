package middleware

//func TestMiddleware_ValidToken(t *testing.T) {
//	maker, _ := pisato.NewPasetoMaker(config.PrivateKey)
//	require.NotEmpty(t, maker)
//
//	token, _ := maker.CreateToken("fahrian a", 4*time.Minute)
//	require.NotEmpty(t, token)
//
//	app := fiber.New()
//	req := &fasthttp.RequestCtx{}
//	ctx := app.AcquireCtx(req)
//	defer app.ReleaseCtx(ctx)
//
//	require.NotNil(t, ctx)
//
//	cookie := &fiber.Cookie{
//		Name:     "token",
//		Value:    token,
//		Expires:  time.Now().Add(4 * time.Hour),
//		Secure:   false,
//		HTTPOnly: true,
//	}
//
//	ctx.Cookie(cookie)
//
//	require.NotEmpty(t, ctx.Cookies("token"))
//	err := Middleware(ctx)
//	require.NoError(t, err)
//
//	assert.Nil(t, err, "Expected error to be nil")
//	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode(), "Expected status code 200")
//}
