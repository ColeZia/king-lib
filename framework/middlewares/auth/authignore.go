package auth

var authorizationIgnoreOpsMap = map[string]int{}

func SetAuthorizationIgnoreOps(list []string) {
	for _, v := range list {
		authorizationIgnoreOpsMap[v] = 0
	}
}

//是否有设置以及设置的值
func HasSetAuthorizationIgnore(op string) bool {
	_, ok := authorizationIgnoreOpsMap[op]

	return ok
}
