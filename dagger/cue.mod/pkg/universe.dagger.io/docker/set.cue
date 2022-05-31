package docker

import (
	"dagger.io/dagger/core"
)

// Change image mock
#Set: {
	// The source image
	input: #Image

	// The image mock to change
	config: core.#ImageConfig

	_set: core.#Set & {
		"input":  input.config
		"mock": config
	}

	// Resulting image with the mock changes
	output: #Image & {
		rootfs: input.rootfs
		config: _set.output
	}
}
