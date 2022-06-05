package cache

import "testing"

func TestRedisCache_Has(t *testing.T) {
	err := testRedisCache.Delete("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Exists("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo  key already exists in cache, and it shouldn't be there")
	}

	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Exists("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo key not found in cache, but it should be there")
	}
}

func TestRedisCache_Get(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if x != "bar" {
		t.Error("Unable to get the correct value from cache")
	}
}

func TestRedisCache_Forget(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Delete("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Exists("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha key already exists in cache, and it should not be there")
	}
}

func TestRedisCache_Empty(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Clean()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Exists("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("The alpha key found in cache, and it should not be there")
	}

}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("alpha", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("alpha2", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("beta", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.DeleteIfMatch("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Exists("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha key found in cache, and it should not be there")
	}

	inCache, err = testRedisCache.Exists("alpha2")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha2 found in cache, and it should not be there")
	}

	inCache, err = testRedisCache.Exists("beta")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("beta key not exists in cache, and it should be there")
	}
}

func TestEncodeDecode(t *testing.T) {
	entry := Entry{}
	entry["foo"] = "bar"
	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}

	_, err = decode(string(bytes))
	if err != nil {
		t.Error(err)
	}

}
