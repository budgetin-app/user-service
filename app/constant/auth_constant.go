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
func GetUserRoleNameByID(id uint) string {
	if int(id) <= len(roleNames) {
		return roleNames[id-1]
	} else {
		panic(fmt.Sprintf("user role is not within the range, id '%d'", id))
	}
}

func GetUserRoles() map[uint]string {
	list := make(map[uint]string, len(roleIds))
	for _, id := range roleIds {
		list[id] = GetUserRoleNameByID(id)
	}
	return list
}
