package valve

var local bool
var C Client

func init() {
	local = true
	if local {
		var err error
		C, err = NewClient(local)
		if err != nil {
			panic(err)
		}
	}
}
