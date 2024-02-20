package constant

import "fmt"

const (
	UserRoleID uint = iota + 1
	AdminRoleID
	// .. specifiy other roles here
)

// roleNames should be updated when the constants changed,
// the ids and names order matter!
var roleIds = []uint{UserRoleID, AdminRoleID}
var roleNames = []string{"User", "Admin"}

// GetUserRoleNameByID get the name of the user role by it's id
func GetUserRoleNameByID(id uint) (*string, error) {
	if int(id) <= len(roleNames) {
		return &roleNames[id-1], nil
	} else {
		return nil, fmt.Errorf("user role with id %d not found", id)
	}
}

func GetUserRoles() map[uint]string {
	list := make(map[uint]string, len(roleIds))
	for _, id := range roleIds {
		if role, err := GetUserRoleNameByID(id); err == nil {
			list[id] = *role
		}
	}
	return list
}
