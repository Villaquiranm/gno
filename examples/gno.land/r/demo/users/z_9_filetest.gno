package main

import (
	"std"

	"gno.land/p/demo/testutils"
	"gno.land/r/demo/users"
)

const admin = std.Address("g1manfred47kzduec920z88wfr64ylksmdcedlf5")

func main() {
	caller := std.OriginCaller() // main
	test2 := testutils.TestAddress("test2")
	// as admin, invite gnouser and test2
	std.TestSetOriginCaller(admin)
	users.Invite(caller.String() + "\n" + test2.String())
	// register as caller
	std.TestSetOriginCaller(caller)
	users.Register(admin, "gnouser", "my profile")
	// register as test2
	std.TestSetOriginCaller(test2)
	users.Register(admin, "test222", "my profile 2")
	println("done")
}

// Output:
// done
