package main

import gen "niaoshenhao.com/gen/modelGen"

func main() {
	//model path
	goModelPath := "/Users/cooerson/Documents/go/src/niaoshenhao.com/ibis/models"
	iosModelPath := "/Users/cooerson/Documents/ios/Ibis-ios/Ibiscoin/Ibiscoin/Model"

	//forms
	gen.GoGenForm(goModelPath, "forms")
	gen.OCGenForm(iosModelPath, "Forms")
	//returns
	gen.GoGenReturn(goModelPath, "returns")
	gen.OCGenReturn(iosModelPath, "Returns")

	//user
	gen.GoGen(goModelPath, "user")
	gen.OCGen(iosModelPath, "User")
	//role
	gen.GoGen(goModelPath, "role")
	gen.OCGen(iosModelPath, "Role")
	//skill
	gen.GoGen(goModelPath, "skill")
	gen.OCGen(iosModelPath, "Skill")
	//coin
	gen.GoGen(goModelPath, "coin")
	gen.OCGen(iosModelPath, "Coin")
	//msg
	gen.GoGen(goModelPath, "msg")
	gen.OCGen(iosModelPath, "Msg")

	// gen.GoGen(goModelPath, "love")
	// gen.OCGen(iosModelPath, "Love")

	// gen.GoGen(goModelPath, "note")
	// gen.OCGen(iosModelPath, "Note")
}
