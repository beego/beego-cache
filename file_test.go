// Copyright 2021 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileCacheInit(t *testing.T) {
	fc := &FileCache{}
	FileCacheWithCachePath("////aaa")(fc)
	err := fc.Init()
	assert.NotNil(t, err)
	FileCacheWithCachePath(getTestCacheFilePath())(fc)
	err = fc.Init()
	assert.Nil(t, err)
}

func TestFileGetContents(t *testing.T) {
	_, err := FileGetContents("/bin/aaa")
	assert.NotNil(t, err)
	fn := filepath.Join(os.TempDir(), "fileCache.txt")
	f, err := os.Create(fn)
	assert.Nil(t, err)
	_, err = f.WriteString("text")
	assert.Nil(t, err)
	data, err := FileGetContents(fn)
	assert.Nil(t, err)
	assert.Equal(t, "text", string(data))
}

func TestGobEncodeDecode(t *testing.T) {
	_, err := GobEncode(func() {
		fmt.Print("test func")
	})
	assert.NotNil(t, err)
	data, err := GobEncode(&FileCacheItem{
		Data: "hello",
	})
	assert.Nil(t, err)
	err = GobDecode([]byte("wrong data"), &FileCacheItem{})
	assert.NotNil(t, err)
	dci := &FileCacheItem{}
	err = GobDecode(data, dci)
	assert.Nil(t, err)
	assert.Equal(t, "hello", dci.Data)
}

func TestFileCacheDelete(t *testing.T) {
	fc, err := NewFileCache()
	assert.Nil(t, err)
	err = fc.Delete(context.Background(), "my-key")
	assert.Nil(t, err)
}

func getTestCacheFilePath() string {
	return filepath.Join(os.TempDir(), "test", "file.txt")
}
