//go:build mage
// +build mage

package main

// Creates the binary in the current directory.  It will overwrite any existing
// binary.
func Build() {
	print("building!")
	/*	set GOOS=linux
		set GOARCH=arm
		set GOARM=7
		set CGO_ENABLED="1"
		set GITHASH=""
		set TIMESTAMP=""


		set TIMESTAMP=%TIMESTAMP: =_%
	*/
	ldf = "-X 'main.BubblesnetBuildNumberString=201' -X 'main.BubblesnetVersionMajorString=2' -X 'main.BubblesnetVersionMinorString=1' -X 'main.BubblesnetVersionPatchString=1'  -X 'main.BubblesnetGitHash=$GITHASH' -X main.BubblesnetBuildTimestamp=$TIMESTAMP"

	return run("go", "build", "-tags", "make", "--ldflags="+ldf, "-o", "build", "./...")

}

// Sends the binary to the server.
func Deploy() error {
	return nil
}
