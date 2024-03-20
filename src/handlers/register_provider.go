package handlers

//var _ Handler = new(registerProviderHandler)
//
//type registerProviderHandler struct{}
//
//func (o registerProviderHandler) Handle(_ context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
//	log.Println("enter registerProviderHandler.Handle()")
//
//	userToken := req.Headers["authtoken"]
//	if userToken == "" {
//		return responders.Error{
//			Code: http.StatusUnauthorized,
//			Err:  errors.New("AuthToken required"),
//			Req:  &req,
//		}.Respond()
//	}
//
//	err := utils.ValidateToken(req.Headers[utils.AuthHeader])
//	authToken, err := utils.SecretsManagerToken()
//	if err != nil {
//		return responders.Error{
//			Code: http.StatusInternalServerError,
//			Err:  fmt.Errorf("failed while getting secret - %w", err),
//		}.Respond()
//	}
//
//	if userToken != authToken {
//		return responders.Error{
//			Code: http.StatusUnauthorized,
//			Err:  errors.New("AuthToken not valid"),
//			Req:  &req,
//		}.Respond()
//	}
//
//	r := responders.RegisterProvider{}
//	if req.IsBase64Encoded {
//		var body []byte
//		body, err = base64.StdEncoding.DecodeString(req.Body)
//		if err != nil {
//			return responders.Error{
//				Code: http.StatusInternalServerError,
//				Err:  fmt.Errorf("failed decoding base64 body - %w", err),
//				Req:  &req,
//			}.Respond()
//		}
//
//		err = json.Unmarshal(body, &r)
//	} else {
//		err = json.Unmarshal([]byte(req.Body), &r)
//	}
//	if err != nil {
//		return responders.Error{
//			Code: http.StatusInternalServerError,
//			Err:  fmt.Errorf("failed base64 decoding request body - %w", err),
//			Req:  &req,
//		}.Respond()
//	}
//
//	return r.Respond()
//}
//
//func newRegisterProviderHandler() Handler {
//	return wellKnownHandler{}
//}
