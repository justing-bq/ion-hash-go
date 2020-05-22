/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package ionhash

import "fmt"

// An InvalidOperationError is returned when a method call is invalid for the struct's current state.
type InvalidOperationError struct {
	structName string
	methodName string
	errorMessage string
}

func (e *InvalidOperationError) Error() string {
	if e.errorMessage != "" {
		return fmt.Sprintf(`ionhash: %v.%v: %v`, e.structName, e.methodName, e.errorMessage)
	} else {
		return fmt.Sprintf(`ionhash: invalid operation error in %v.%v`, e.structName, e.methodName)
	}
}

// InvalidArgumentError is returned when one of the arguments given to a function was not valid.
type InvalidArgumentError struct {
	argumentName string
	argumentValue interface{}
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf(`ionhash: invalid value: "%v" specified for argument: %s`, e.argumentValue, e.argumentName)
}
