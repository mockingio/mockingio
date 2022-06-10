package main

import (
    "dagger.io/dagger"
		"dagger.io/dagger/core"

    "universe.dagger.io/go"

		"github.com/smockyio/dagger/ci/golangci"
)

dagger.#Plan & {
    client: filesystem: ".": read: contents: dagger.#FS
		client: filesystem: "./bin": write: contents: actions.build."go".output

    actions: {
				_source: client.filesystem.".".read.contents

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
							os: *client.platform.os | "linux"
							arch: client.platform.arch
							ldflags: "-s -w -X github.com/smockyio/smocky/backend/cmd/cli.buildVersion=1.0.1"
							env: depends_unit: "\(test.unit.exit)"
					}

					docker: core.#Dockerfile & {
						source: _source
						dockerfile: path: "Dockerfile"
					}
        }
    }

}