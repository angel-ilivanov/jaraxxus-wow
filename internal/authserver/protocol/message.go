package protocol

type Request interface {
	isRequest()
}

type Response interface {
	isResponse()
}
