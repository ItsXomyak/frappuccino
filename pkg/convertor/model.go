package convertor

type Massage struct {
	Status  int    `json:"status"`
	Messege string `json:"messege"`
}

func Wrap(status int, message string) *Massage {
	custom := &Massage{
		Status:  status,
		Messege: message,
	}

	return custom
}
