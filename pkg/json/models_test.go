package json

import (
	"github.com/AkronimBlack/process-manager/shared"
	"reflect"
	"testing"
)

func TestArgs_GetMap(t *testing.T) {
	testForMap := map[string]interface{}{
		"key_1_sub_1": "value",
		"key_2_sub_2": map[string]string{
			"test": "test",
		},
		"key_3_sub_3": []string{
			"test_1", "test_2", "test_3",
		},
	}
	args := Args{
		"key_1": "test_1",
		"key_2": testForMap,
	}
	subMap := args.GetMap("key_2")
	if !reflect.DeepEqual(testForMap, subMap) {
		t.Errorf("maps not the same\n expected: %s\n got: %s", shared.ToJsonPrettyString(testForMap), shared.ToJsonPrettyString(subMap))
	}
}

func TestArgs_GetMapReturnsDefault(t *testing.T) {
	testForMap := map[string]interface{}{
		"key_1_sub_1": "value",
		"key_2_sub_2": map[string]string{
			"test": "test",
		},
		"key_3_sub_3": []string{
			"test_1", "test_2", "test_3",
		},
	}
	args := Args{
		"key_1": "test_1",
		"key_2": testForMap,
	}
	subMap := args.GetMap("key_3")
	if len(subMap) != 0 {
		t.Errorf("maps not the same\n expected: %s\n got: %s", shared.ToJsonPrettyString(testForMap), shared.ToJsonPrettyString(subMap))
	}

	defaultMap := map[string]interface{}{"test": "test"}
	subMap = args.GetMap("key_3", defaultMap)
	if !reflect.DeepEqual(defaultMap, subMap) {
		t.Errorf("maps not the same\n expected: %s\n got: %s", shared.ToJsonPrettyString(defaultMap), shared.ToJsonPrettyString(subMap))
	}
}

func TestArgs_Bind(t *testing.T) {
	args := Args{
		"key_1": "value",
		"key_2": map[string]string{
			"test": "test",
		},
		"key_3": []string{
			"test_1", "test_2", "test_3",
		},
	}

	type test struct {
		Key1 string            `json:"key_1"`
		Key2 map[string]string `json:"key_2"`
		Key3 []string          `json:"key_3"`
	}
	v := test{}

	err := args.Bind(&v)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(shared.ToJsonPrettyString(v))
}
