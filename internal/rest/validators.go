package rest

import "fmt"

func validateRegister(req *RegisterRequestBody) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func validateLogin(req *LoginRequestBody) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if req.AppId == 0 {
		return fmt.Errorf("app_id is required")
	}

	return nil
}