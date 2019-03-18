// date: 2019-03-18
package broker

type ExternalAuthentication struct {
}

func (a *ExternalAuthentication) Mechanism() string {
	return "EXTERNAL"
}

func (a *ExternalAuthentication) Response() string {
	return ""

}
