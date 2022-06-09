package main

import (
    "dagger.io/dagger"

		"github.com/smockyio/dagger/ci/golangci"

    "universe.dagger.io/go"
)

dagger.#Plan & {
    client: filesystem: ".": read: contents: dagger.#FS

    actions: {
				_source: client.filesystem.".".read.contents
				_root_source: client.filesystem["."].read.contents

				test: {
					unit: go.#Test & {
						source:  _source
						package: "./..."
						command: flags: "-race": true
					}
				}

				lint: {
						go: golangci.#Lint & {
						source:  _source
						version: "1.45"
					}
				}

        build: {
        	"go": go.#Build & {
							source: client.filesystem.".".read.contents
							env: {
								CGO_ENABLED: "0"
							}
							package: "."
							env: HACK: "\(test.unit.success)"
					}
        }
    }

}