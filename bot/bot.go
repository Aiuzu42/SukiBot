package bot

import (
	"strings"

	"github.com/aiuzu42/SukiBot/config"
)

var (
	owners = []string{}
	admins = []string{}
)

const (
	prefix = "suki-"
	pLen   = 5
)

func LoadRoles() {
	owners = config.Config.Users.Owners
	admins = config.Config.Users.Admins
}

// IsOwner returns true if userID matches an ID of the owners group.
func IsOwner(userID string) bool {
	return findIfExists(userID, owners)
}

// IsAdmin returns true if any of the roles is part of the admins group or if the userID matches an ID of the owners group.
func IsAdmin(roles []string, userID string) bool {
	if IsOwner(userID) {
		return true
	}
	return arrayFindIfExists(roles, admins)
}

func findIfExists(a string, b []string) bool {
	for _, eb := range b {
		if a == eb {
			return true
		}
	}
	return false
}

func arrayFindIfExists(a []string, b []string) bool {
	for _, ea := range a {
		for _, eb := range b {
			if ea == eb {
				return true
			}
		}
	}
	return false
}

func argumentsHandler(st string) (string, string) {
	arg := ""
	msg := ""
	args := strings.Split(st, " ")
	n := len(args)
	if n == 2 {
		arg = args[1]
	} else if n > 2 {
		arg = args[1]
		msg = strings.Join(args[2:], " ")
	}
	return arg, msg
}
