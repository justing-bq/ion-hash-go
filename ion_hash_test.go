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
	"bytes"
	"fmt"
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

	for reader.Next() {
		assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

		reader.Next() //reads the ion value eg. 5

		// Create binary writer to write reader Ion values and create binary
		buf := bytes.Buffer{}
		writer := ion.NewBinaryWriter(&buf)
		writeValue(t, reader, writer)
		assert.NoError(t, writer.Finish(), "Something went wrong writing Ion value to binary writer.")

		// Create reader with the binary to create hash reader
		hr, err := NewHashReader(ion.NewReaderBytes(buf.Bytes()), newIdentityHasherProvider())
		assert.NoError(t, err, "Something went wrong creating hash reader.")

		hr.Next()
		hr.Next()
		fmt.Println(hr.Sum(nil))

		reader.Next()
		fieldName := reader.FieldName()
		if fieldName == "expect" {
			assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			reader.Next()

			fieldName = reader.FieldName()
			if fieldName == "identity" {
				assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

				for reader.Next() {
					annotations := reader.Annotations()

					if len(annotations) > 0 {
						if annotations[0] == "update" {
							fmt.Println(reader.ByteValue())
						} else if annotations[0] == "digest" {
							fmt.Println(reader.ByteValue())
						}
					}
				}
				assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
			}

			assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
		}

		assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
	}
}

// Get value from reader and write to writer.
func writeValue(t *testing.T, reader ion.Reader, writer ion.Writer) {
	ionType := reader.Type()
	switch ionType {
	case ion.BoolType:
		boolValue, err := reader.BoolValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteBool(boolValue), "Something went wrong executing writer.WriteBool(boolValue)")
	case ion.BlobType:
		byteValue, err := reader.ByteValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteBlob(byteValue), "Something went wrong executing writer.WriteBlob(byteValue)")
	case ion.ClobType:
		byteValue, err := reader.ByteValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteClob(byteValue), "Something went wrong executing writer.WriteClob(byteValue)")
	case ion.DecimalType:
		decimalValue, err := reader.DecimalValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteDecimal(decimalValue), "Something went wrong executing writer.WriteDecimal(decimalValue)")
	case ion.FloatType:
		floatValue, err := reader.FloatValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteFloat(floatValue), "Something went wrong executing writer.WriteFloat(floatValue)")
	case ion.IntType:
		intSize, err := reader.IntSize()
		require.NoError(t, err)

		switch intSize {
		case ion.Int32:
			intValue, err := reader.IntValue()
			require.NoError(t, err)

			assert.NoError(t, writer.WriteInt(int64(intValue)), "Something went wrong executing writer.WriteInt(int64(intValue))")
		case ion.Int64:
			intValue, err := reader.Int64Value()
			require.NoError(t, err)

			assert.NoError(t, writer.WriteInt(intValue), "Something went wrong executing writer.WriteInt(intValue)")
		case ion.Uint64:
			intValue, err := reader.Uint64Value()
			require.NoError(t, err)

			assert.NoError(t, writer.WriteUint(intValue), "Something went wrong executing writer.WriteUint(intValue)")
		case ion.BigInt:
			intValue, err := reader.BigIntValue()
			require.NoError(t, err)

			assert.NoError(t, writer.WriteBigInt(intValue), "Something went wrong executing writer.WriteBigInt(intValue)")
		default:
			t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
		}
	case ion.StringType:
		stringValue, err := reader.StringValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteString(stringValue), "Something went wrong executing writer.WriteString(stringValue)")
	case ion.SymbolType:
		stringValue, err := reader.StringValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteSymbol(stringValue), "Something went wrong executing writer.WriteSymbol(stringValue)")
	case ion.TimestampType:
		timeValue, err := reader.TimeValue()
		require.NoError(t, err)

		assert.NoError(t, writer.WriteTimestamp(timeValue), "Something went wrong executing writer.WriteTimestamp(timeValue)")
	case ion.SexpType:
		require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

		require.NoError(t, writer.BeginSexp(), "Something went wrong executing writer.BeginSexp()")

		writeValue(t, reader, writer)

		require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

		require.NoError(t, writer.EndSexp(), "Something went wrong executing writer.EndSexp()")
	case ion.ListType:
		require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

		require.NoError(t, writer.BeginList(), "Something went wrong executing writer.BeginList()")

		writeValue(t, reader, writer)

		require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

		require.NoError(t, writer.EndList(), "Something went wrong executing writer.EndList()")

	case ion.StructType:
		require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

		require.NoError(t, writer.BeginStruct(), "Something went wrong executing writer.BeginStruct()")

		writeValue(t, reader, writer)

		require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

		require.NoError(t, writer.EndStruct(), "Something went wrong executing writer.EndStruct()")
	default:
		t.Fatal(InvalidIonTypeError{ionType})
	}
}
