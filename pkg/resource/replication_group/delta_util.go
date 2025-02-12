// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package replication_group

import (
	"regexp"
	"strings"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"

	"github.com/aws-controllers-k8s/elasticache-controller/pkg/common"
)

// remove non-meaningful differences from delta
func filterDelta(
	delta *ackcompare.Delta,
	desired *resource,
	latest *resource,
) {

	if delta.DifferentAt("Spec.EngineVersion") {
		if desired.ko.Spec.EngineVersion != nil && latest.ko.Spec.EngineVersion != nil {
			if engineVersionsMatch(*desired.ko.Spec.EngineVersion, *latest.ko.Spec.EngineVersion) {
				common.RemoveFromDelta(delta, "Spec.EngineVersion")
			}
		}
		// TODO: handle the case of a nil difference (especially when desired EV is nil)
	}
}

// returns true if desired and latest engine versions match and false otherwise
// precondition: both desiredEV and latestEV are non-nil
// this handles the case where only the major EV is specified, e.g. "6.x" (or similar), but the latest
//   version shows the minor version, e.g. "6.0.5"
func engineVersionsMatch(
	desiredEV string,
	latestEV string,
) bool {
	if desiredEV == latestEV {
		return true
	}

	// if the last character of desiredEV is "x", only check for a major version match
	last := len(desiredEV) - 1
	if desiredEV[last:] == "x" {
		// cut off the "x" and replace all occurrences of '.' with '\.' (as '.' is a special regex character)
		desired := strings.Replace(desiredEV[:last], ".", "\\.", -1)
		r, _ := regexp.Compile(desired + ".*")
		return r.MatchString(latestEV)
	}

	return false
}
