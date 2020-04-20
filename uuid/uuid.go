package uuid

import (
	"log"
	"os/exec"
)

func UUID() string {
	Out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	uuid := string(Out)

	uuidlen := len(uuid)

	uuid = uuid[:uuidlen-1]

	return uuid
	// fmt.Printf("%s", out)
}
