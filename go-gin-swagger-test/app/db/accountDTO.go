package db

type AccountDTO struct {
	Name string `json:"name"`
}

func (a AccountDTO) IsNameValid() error {
	if len(a.Name) == 0 {
		return ErrNoNameValid
	}

	return nil
}
