package main

import (
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/implementations"
	"reflect"
	"testing"
)

// benchmarks result for ordinary PC (server and client on the same PC)
//Method		Map		Flat	Json	Gob		Rejson
//Set			55657	49625	45790	54854	52233
//Get			57728	49081	48700	73354	50633
//SetOneField	44168	43306	99398	124524	48715
//GetOneField	43283	43784	51440	73798	46301
//Del			42433	42401	42009	42361	41566

//$ go test -bench=.
//Failed to find configuration config-custom.json
//Failed to find configuration config-custom.json
//goos: linux
//goarch: amd64
//pkg: github.com/PROger4ever-Golang/Redis-serialization-benchmarks
//BenchmarkImplementations/Map/Set-4     	   30000	     55657 ns/op
//BenchmarkImplementations/Map/Get-4     	   30000	     57728 ns/op
//BenchmarkImplementations/Map/SetOneField-4         	   30000	     44168 ns/op
//BenchmarkImplementations/Map/GetOneField-4         	   30000	     43283 ns/op
//BenchmarkImplementations/Map/Del-4                 	   30000	     42433 ns/op
//BenchmarkImplementations/Flat/Set-4                	   30000	     49625 ns/op
//BenchmarkImplementations/Flat/Get-4                	   30000	     49081 ns/op
//BenchmarkImplementations/Flat/SetOneField-4        	   30000	     43306 ns/op
//BenchmarkImplementations/Flat/GetOneField-4        	   30000	     43784 ns/op
//BenchmarkImplementations/Flat/Del-4                	   30000	     42401 ns/op
//BenchmarkImplementations/Json/Set-4                	   30000	     45790 ns/op
//BenchmarkImplementations/Json/Get-4                	   30000	     48700 ns/op
//BenchmarkImplementations/Json/SetOneField-4        	   20000	     99398 ns/op
//BenchmarkImplementations/Json/GetOneField-4        	   30000	     51440 ns/op
//BenchmarkImplementations/Json/Del-4                	   30000	     42009 ns/op
//BenchmarkImplementations/Gob/Set-4                 	   30000	     54854 ns/op
//BenchmarkImplementations/Gob/Get-4                 	   20000	     73354 ns/op
//BenchmarkImplementations/Gob/SetOneField-4         	   10000	    124524 ns/op
//BenchmarkImplementations/Gob/GetOneField-4         	   20000	     73798 ns/op
//BenchmarkImplementations/Gob/Del-4                 	   30000	     42361 ns/op
//BenchmarkImplementations/Rejson/Set-4              	   30000	     52233 ns/op
//BenchmarkImplementations/Rejson/Get-4              	   30000	     50633 ns/op
//BenchmarkImplementations/Rejson/SetOneField-4      	   30000	     48715 ns/op
//BenchmarkImplementations/Rejson/GetOneField-4      	   30000	     46301 ns/op
//BenchmarkImplementations/Rejson/Del-4              	   30000	     41566 ns/op
//PASS
//ok  	github.com/PROger4ever-Golang/Redis-serialization-benchmarks	48.786s

func TestImplementations(t *testing.T) {
	app, err := NewApplication()
	if err != nil {
		t.Fatalf(ERROR_FORMAT, "main", "configuring app", err)
		return
	}
	defer app.Finalize()

	for name, imp := range app.Implementations {
		imp.Del()
		srcObject := &implementations.SerializationObject{
			String1: "String1", String2: "String2", String3: "String3", String4: "String4", String5: "String5",
			FieldX: "FieldX",
		}
		t.Run(name, func(t *testing.T) {
			err = imp.Set(srcObject)
			if err != nil {
				t.Fatalf(ERROR_FORMAT, name, "setting object", err)
			}

			object1, err := imp.Get()
			if err != nil {
				t.Fatalf(ERROR_FORMAT, name, "getting object", err)
			}

			isEqual := reflect.DeepEqual(srcObject, object1)
			if !isEqual {
				t.Fatalf(ERROR_FORMAT, name, "comparing src srcObject and object1", "they aren't equal")
			}

			srcObject.FieldX = "FieldX-changed"
			err = imp.SetOneField("FieldX", srcObject.FieldX)
			if err != nil {
				t.Fatalf(ERROR_FORMAT, name, "setting one field", err)
			}

			field, err := imp.GetOneField("FieldX")
			if err != nil {
				t.Fatalf(ERROR_FORMAT, name, "getting one field", err)
			}

			if field != srcObject.FieldX {
				t.Fatalf(ERROR_FORMAT, name, "comparing src field and actual field", "they aren't equal")
			}

			object2, err := imp.Get()
			if err != nil {
				t.Fatalf(ERROR_FORMAT, name, "getting object", err)
			}

			isEqual = reflect.DeepEqual(srcObject, object2)
			if !isEqual {
				t.Fatalf(ERROR_FORMAT, name, "comparing src srcObject and object2", "they aren't equal")
			}

			err = imp.Del()
			if err != nil {
				t.Fatalf(ERROR_FORMAT, name, "deleting object", err)
			}
		})
	}
}

func BenchmarkImplementations(b *testing.B) {
	app, err := NewApplication()
	if err != nil {
		b.Fatalf(ERROR_FORMAT, "main", "configuring app", err)
		return
	}
	defer app.Finalize()

	for name, imp := range app.Implementations {
		imp.Del()
		srcObject := &implementations.SerializationObject{
			String1: "String1", String2: "String2", String3: "String3", String4: "String4", String5: "String5",
			FieldX: "FieldX",
		}
		b.Run(name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.Run("Set", func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						err = imp.Set(srcObject)
						if err != nil {
							b.Fatalf(ERROR_FORMAT, name, "setting object", err)
						}
					}
				})

				b.Run("Get", func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						_, err := imp.Get()
						if err != nil {
							b.Fatalf(ERROR_FORMAT, name, "getting object", err)
						}
					}
				})

				srcObject.FieldX = "FieldX-changed"
				b.Run("SetOneField", func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						err = imp.SetOneField("FieldX", srcObject.FieldX)
						if err != nil {
							b.Fatalf(ERROR_FORMAT, name, "setting one field", err)
						}
					}
				})

				b.Run("GetOneField", func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						_, err := imp.GetOneField("FieldX")
						if err != nil {
							b.Fatalf(ERROR_FORMAT, name, "getting one field", err)
						}
					}
				})

				b.Run("Del", func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						err = imp.Del()
						if err != nil {
							b.Fatalf(ERROR_FORMAT, name, "deleting object", err)
						}
					}
				})
			}
		})
	}
}
