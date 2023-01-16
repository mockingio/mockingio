package mock

import "github.com/mockingio/mockingio/engine/mock"

type Method struct {
	builder *Builder
}

// Response is used to define what will be responded to the request.
func (m *Method) Response(status int, body string, headers ...Headers) *Response {
	if len(headers) == 0 {
		headers = append(headers, map[string]string{})
	}

	m.builder.response = &mock.Response{
		Body:    body,
		Status:  status,
		Headers: headers[0],
	}

	resp := &Response{
		builder: m.builder,
	}

	return resp
}
