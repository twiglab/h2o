package wecom

import "cmp"

type MessageToken struct {
	ChatID   string
	ChatType string
	UserID   string
}

type FixUserGroup struct {
	Users  []string
	ChatID string
}

func (a FixUserGroup) Auth(t MessageToken) (bool, error) {
	if t.ChatType != "group" {
		return false, nil
	}

	if cmp.Compare(a.ChatID, t.ChatID) != 0 {
		return false, nil
	}

	for _, u := range a.Users {
		if cmp.Compare(u, t.UserID) == 0 {
			return true, nil
		}
	}

	return false, nil
}
