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
