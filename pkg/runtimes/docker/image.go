/*
Copyright © 2020-2022 The k3d Author(s)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

// GetImages returns a list of images present in the runtime
func (d Docker) GetImages(ctx context.Context) ([]string, error) {
	// create docker client
	docker, err := GetDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer docker.Close()

	imageSummary, err := docker.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("docker failed to list images: %w", err)
	}

	var images []string
	for _, image := range imageSummary {
		images = append(images, image.RepoTags...)
		images = append(images, image.RepoDigests...)
		images = append(images, repoTagDigests(image)...)
	}

	return images, nil
}

func repoTagDigests(image types.ImageSummary) []string {
	if len(image.RepoTags) == 0 {
		return []string{}
	}

	base := strings.Split(image.RepoTags[0], ":")[0]
	tags := []string{}
	for _, tag := range image.RepoTags {
		tags = append(tags, strings.Split(tag, ":")[1])
	}

	digests := []string{}
	for _, digest := range image.RepoDigests {
		digests = append(digests, strings.Split(digest, "@")[1])
	}

	tagDigests := []string{}
	for _, tag := range tags {
		for _, digest := range digests {
			tagDigests = append(tagDigests, fmt.Sprintf("%s:%s@%s", base, tag, digest))
		}
	}
	return tagDigests
}
