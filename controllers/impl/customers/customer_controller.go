package customers

import "net/http"

type CustomerController struct {
}

func NewCustomerController() *CustomerController {
	return &CustomerController{}
}

func (c *CustomerController) GetCustomers(w http.ResponseWriter, r *http.Request) {}
