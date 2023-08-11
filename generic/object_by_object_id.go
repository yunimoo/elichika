package generic

import (
	"elichika/common"

	"fmt"

	"encoding/json"
)

// Common pattern of [object_1_id, object_1, object_2_id, object_2...]
// if null give an empty array

type ObjectByObjectIDReadInterface interface { // only unmarshal
	SetID(int64)
}

type ObjectByObjectIDRead[T ObjectByObjectIDReadInterface] struct {
	Objects []T
}

type ObjectByObjectIDWriteInterface interface { // only marshal
	ID() int64
}

type ObjectByObjectIDWrite[T ObjectByObjectIDWriteInterface] struct {
	Objects []T
}

func (oboid ObjectByObjectIDWrite[T]) MarshalJSON() ([]byte, error) {
	arr := []any{}
	for _, object := range oboid.Objects {
		arr = append(arr, object.ID())
		arr = append(arr, object)
	}
	bytes, err := json.Marshal(&arr)
	common.CheckErr(err)
	fmt.Println("good")
	return bytes, nil
}

func (oboid *ObjectByObjectIDRead[T]) UnmarshalJSON(data []byte) error {
	oboid.Objects = []T{}
	if string(data) == "null" {
		return nil
	}
	arr := []any{}
	err := json.Unmarshal(data, &arr)
	common.CheckErr(err)
	for i, _ := range arr {
		if i%2 == 1 {
			continue
		}
		bytes, err := json.Marshal(arr[i+1])
		common.CheckErr(err)
		obj := new(T)
		err = json.Unmarshal(bytes, &obj)
		(*obj).SetID(int64(arr[i].(float64)))
		common.CheckErr(err)
		oboid.Objects = append(oboid.Objects, *obj)
	}
	return nil
}
