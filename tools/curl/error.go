package curl

import "fmt"

type CurlError struct {
	name    string      // Task struct Name
	code    int         // Task struct Code
	message interface{} // Error message
}

func (this CurlError) Error() string {
	name := fmt.Sprintf("Name  : %v\n", this.name)
	code := fmt.Sprintf("Code  : %v\n", this.code)
	msg := fmt.Sprintf("Error : %v", this.message)
	return "\n" + name + code + msg
}
