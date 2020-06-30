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

import (
	"io/ioutil"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIonHash(t *testing.T) {
	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	reader := ion.NewReaderBytes(file)

	var testCases []interface{}

	for reader.Next() {
		err := reader.StepIn()
		assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

		reader.Next()

		testName := "unknown"

		reader.Next()
		fieldName := reader.FieldName()
		if fieldName == "expect" {
			err = reader.StepIn()
			assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

			reader.Next()

			fieldName = reader.FieldName()
			if fieldName == "identity" {
				err = reader.StepIn()
				assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

				for reader.Next() {
					annotations := reader.Annotations()

					if len(annotations) > 0 {
						if annotations[0] == "update" {
							ionVal := IonValue(t, reader)
							testCases = append(testCases, ionVal)
						} else if annotations[0] == "digest" {
							digestVal := IonValue(t, reader)
							if digestVal == nil {

							}
						}
					}
				}
				err = reader.StepOut()
				assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
			}

			err = reader.StepOut()
			assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
		}

		annotations := reader.Annotations()

		if len(annotations) > 0 {
			testName = annotations[0]
		}

		if testName == "" {

		}

		err = reader.StepOut()
		assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
	}
}

func IonValue(t *testing.T, reader ion.Reader) interface{} {
	var ionValue interface{}

	ionType := reader.Type()
	switch ionType {
	case ion.BoolType:
		boolValue, err := reader.BoolValue()
		require.NoError(t, err)

		ionValue = boolValue
	case ion.BlobType, ion.ClobType:
		byteValue, err := reader.ByteValue()
		require.NoError(t, err)

		ionValue = byteValue
	case ion.DecimalType:
		decimalValue, err := reader.DecimalValue()
		require.NoError(t, err)

		ionValue = decimalValue
	case ion.FloatType:
		floatValue, err := reader.FloatValue()
		require.NoError(t, err)

		ionValue = floatValue
	case ion.IntType:
		intSize, err := reader.IntSize()
		require.NoError(t, err)

		switch intSize {
		case ion.Int32:
			intValue, err := reader.IntValue()
			require.NoError(t, err)

			ionValue = intValue
		case ion.Int64:
			intValue, err := reader.Int64Value()
			require.NoError(t, err)

			ionValue = intValue
		case ion.Uint64:
			intValue, err := reader.Uint64Value()
			require.NoError(t, err)

			ionValue = intValue
		case ion.BigInt:
			intValue, err := reader.BigIntValue()
			require.NoError(t, err)

			ionValue = intValue
		default:
			t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
		}
	case ion.StringType:
		stringValue, err := reader.StringValue()
		require.NoError(t, err)

		ionValue = stringValue
	case ion.SymbolType:
		stringValue, err := reader.StringValue()
		require.NoError(t, err)

		ionValue = stringValue
	case ion.TimestampType:
		timeValue, err := reader.TimeValue()
		require.NoError(t, err)

		ionValue = timeValue
	case ion.SexpType, ion.ListType, ion.StructType:
		err := reader.StepIn()
		require.NoError(t, err)

		val := IonValue(t, reader)

		err = reader.StepOut()
		require.NoError(t, err)

		return val
	default:
		t.Fatal(InvalidIonTypeError{ionType})
	}

	return ionValue
}
