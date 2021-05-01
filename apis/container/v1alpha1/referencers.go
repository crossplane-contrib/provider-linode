/*
Copyright 2020 The Crossplane Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/reference"
	resource "github.com/crossplane/crossplane-runtime/pkg/resource"
)

// LKEClusterID extracts the ID of an LKECluster.
func LKEClusterID() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		c, ok := mg.(*LKECluster)
		if !ok {
			return ""
		}
		return strconv.Itoa(c.Status.AtProvider.ID)
	}
}
