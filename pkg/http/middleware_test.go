package http

/*

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authType string,
	userID int64,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateAccess(userID, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set(authHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker service.TokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker service.TokenMaker) {
				addAuthorization(t, req, tokenMaker, authTypeBearer, 1, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuth",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker service.TokenMaker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuth",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker service.TokenMaker) {
				addAuthorization(t, req, tokenMaker, "UnsupportedType", 1, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthFormat",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker service.TokenMaker) {
				addAuthorization(t, req, tokenMaker, "", 1, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker service.TokenMaker) {
				addAuthorization(t, req, tokenMaker, authTypeBearer, 1, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"

			server.router.Pre(authMiddleware(server.token))
			server.router.GET(authPath, func(ctx echo.Context) error {
				return ctx.JSON(http.StatusOK, nil)
			})

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.token)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
*/
