package errorPkg

type CustomErrors interface {
	Error() string
	HttpStatusCode() int
}
