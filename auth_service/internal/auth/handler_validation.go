package auth

import "context"

type Validation struct {
	repo Repository
}

func ValidateSignup(ctx context.Context, name, email string) error {
	var val Validation

	err := val.repo.FindEmail(ctx, email)
	if err != nil {
		return err
	}

	err = val.repo.FindSameName(ctx, name)
	if err != nil {
		return err
	}

	return nil
}
