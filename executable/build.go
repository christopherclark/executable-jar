/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package executable

import (
	"fmt"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libjvm"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.BuildResult{}

	m, err := libjvm.NewManifest(context.Application.Path)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to read manifest in %s: %w", context.Application.Path, err)
	}

	if s, ok := m.Get("Main-Class"); ok {
		cp := []string{context.Application.Path}

		if s, ok := m.Get("Class-Path"); ok {
			cp = append(cp, strings.Split(s, " ")...)
		}

		result.Layers = append(result.Layers, NewClassPath(cp))

		command := fmt.Sprintf(`java -cp "${CLASSPATH}" ${JAVA_OPTS} %s`, s)
		result.Processes = append(result.Processes,
			libcnb.Process{Type: "executable-jar", Command: command},
			libcnb.Process{Type: "task", Command: command},
			libcnb.Process{Type: "web", Command: command},
		)
	}

	return result, nil
}
