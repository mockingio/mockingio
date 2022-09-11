package golangci

import (
	"dagger.io/dagger"

	"universe.dagger.io/docker"
	"universe.dagger.io/go"
)

// Lint using golangci-lint
#Lint: {
	// Source code
	source: dagger.#FS

	// golangci-lint version
	version: *"1.49" | string

	// timeout
	timeout: *"5m" | string

	_image: docker.#Pull & {
		source: "golangci/golangci-lint:v\(version)"
	}

	_goImage: go.#Image & {
		"version": "1.18"
	}

	go.#Container & {
		"source": source
		"image": _goImage.output
		input:    _image.output
		command: {
			name: "golangci-lint"
			flags: {
				run:         true
				"-v":        true
				"--timeout": timeout
			}
		}
	}
}
