package handlers

//var _ Handler = new(registerHandler)
//
//type registerHandler struct{}
//
//func (o registerHandler) Handle(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
//	log.Println("enter registerHandler.Handle()")
//	urlParts := strings.Split(strings.TrimLeft(req.RawPath, "/"), "/")
//
//	var r responders.Responder
//	switch {
//	case len(urlParts) != 2:
//		r = responders.Error{
//			Code: http.StatusNotFound,
//			Err:  errors.New("URL path has wrong part count"),
//			Req:  &req,
//		}
//	case urlParts[1] == "provider":
//		h := registerProviderHandler{}
//		return h.Handle(ctx, req)
//	default:
//		r = responders.Error{
//			Code: http.StatusNotFound,
//			Req:  &req,
//		}
//	}
//
//	return r.Respond()
//}
//
//func newRegisterHandler() Handler {
//	return registerHandler{}
//}
