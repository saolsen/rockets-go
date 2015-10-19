NOTE: This is abandoned and a very imcomplete version of rockets in go, open sourcing to free up a private repo, not because it's something anybody will want to see.

# rockets-go

Need to have gopath set up correctly

build dependencies (osx only for now)
* `brew install premake`

* `git submodule init`
* `git submodule update --recursive`

* `cd nanovg`
* `premake4 gmake`
* `cd build`
* `make nanovg`
* `cd ../..`

* `go get .`

* `go run game.go osx.go`
