package main

import (
    "dagger.io/dagger"
		"dagger.io/dagger/core"

		"universe.dagger.io/bash"
		"universe.dagger.io/alpine"
    "universe.dagger.io/go"

		"github.com/mockingio/dagger/ci/golangci"
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

				version: {
					_image: alpine.#Build & {
						packages: bash: _
						packages: curl: _
						packages: git: _
					}

					_revision: bash.#Run & {
						input:   _image.output
						workdir: "/src"
						mounts: source: {
							dest:     "/src"
							contents: _source
						}

						script: contents: #"""
							printf "$(git rev-parse --short HEAD)" > /revision
							"""#
						export: files: "/revision": string
					}

					output: _revision.export.files["/revision"]
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
							ldflags: "-s -w -X github.com/mockingio/mockingio/cmd/version.Revision=\(version.output)"
							env: depends_unit: "\(test.unit.exit)"
					}

					docker: core.#Dockerfile & {
						source: _source
						dockerfile: path: "Dockerfile"
					}
        }
    }

}