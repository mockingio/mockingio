package main

import (
    "dagger.io/dagger"

		"universe.dagger.io/docker"
    "universe.dagger.io/go"
)

dagger.#Plan & {
    client: filesystem: "./../backend": read: contents: dagger.#FS

    actions: {
				_source: client.filesystem."./../backend".read.contents
				_image: docker.#Pull & {
					source: "golangci/golangci-lint:v1.45"
				}
        server_test: go.#Test & {
            source:  client.filesystem."./../backend".read.contents
            package: "./..."
            env: CGO_ENABLED: "0"
        }

        server_lint: go.#Container & {
					source: _source
					input: _image.output
					command: {
						name: "golangci-lint"
						flags: {
							run:         true
							"-v":        true
						}
					}
				}

        build: go.#Build & {
					source: client.filesystem."./../backend".read.contents
					env: CGO_ENABLED: "0"
					package: "./cmd/cli"
					env: HACK: "\(server_test.success)"
        }
    }

}