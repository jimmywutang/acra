// Copyright 2016, Cossack Labs Limited
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
package keystore

import (
	"fmt"
	"github.com/cossacklabs/acra/utils"
	"os"
	"path/filepath"
	"testing"
)

func testGeneral(store *FilesystemKeyStore, t *testing.T) {
	if store.HasZonePrivateKey([]byte("non-existent key")) {
		t.Fatal("Expected false on non-existent key")
	}
	key, err := store.GetZonePrivateKey([]byte("non-existent key"))
	if err == nil {
		t.Fatal("Expected any error")
	}
	if key != nil {
		t.Fatal("Non-expected key")
	}
	id, _, err := store.GenerateZoneKey()
	if err != nil {
		t.Fatal(err)
	}
	if !store.HasZonePrivateKey(id) {
		t.Fatal("Expected true on existed id")
	}
	key, err = store.GetZonePrivateKey(id)
	if err != nil {
		t.Fatal(err)
	}
	if key == nil {
		t.Fatal("Expected private key")
	}
}

func testGeneratingDataEncryptionKeys(store *FilesystemKeyStore, t *testing.T) {
	testId := []byte("test id")
	err := store.GenerateDataEncryptionKeys(testId)
	if err != nil {
		t.Fatal(err)
	}
	exists, err := utils.FileExists(
		store.getFilePath(
			store.getServerDecryptionKeyFilename(testId)))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("Private decryption key doesn't exists")
	}

	exists, err = utils.FileExists(
		fmt.Sprintf("%s.pub", store.getFilePath(
			store.getServerDecryptionKeyFilename(testId))))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("Public decryption key doesn't exists")
	}
}

func testGenerateServerKeys(store *FilesystemKeyStore, t *testing.T) {
	testId := []byte("test id")
	err := store.GenerateServerKeys(testId)
	if err != nil {
		t.Fatal(err)
	}
	expectedPaths := []string{
		store.getServerKeyFilename(testId),
		fmt.Sprintf("%s.pub", store.getServerKeyFilename(testId)),
	}
	for _, name := range expectedPaths {
		absPath := store.getFilePath(name)
		exists, err := utils.FileExists(absPath)
		if err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Fatal(fmt.Sprintf("File <%s> doesn't exists", absPath))
		}
	}
}

func testGenerateProxyKeys(store *FilesystemKeyStore, t *testing.T) {
	testId := []byte("test id")
	err := store.GenerateProxyKeys(testId)
	if err != nil {
		t.Fatal(err)
	}
	expectedPaths := []string{
		store.getProxyKeyFilename(testId),
		fmt.Sprintf("%s.pub", store.getProxyKeyFilename(testId)),
	}
	for _, name := range expectedPaths {
		absPath := store.getFilePath(name)
		exists, err := utils.FileExists(absPath)
		if err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Fatal(fmt.Sprintf("File <%s> doesn't exists", absPath))
		}
	}
}

func testReset(store *FilesystemKeyStore, t *testing.T) {
	testId := []byte("some test id")
	if err := store.GenerateServerKeys(testId); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetServerPrivateKey(testId); err != nil {
		t.Fatal(err)
	}
	store.Reset()
	if err := os.Remove(store.getFilePath(store.getServerKeyFilename(testId))); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(fmt.Sprintf("%s.pub", store.getFilePath(store.getServerKeyFilename(testId)))); err != nil {
		t.Fatal(err)
	}

	if _, err := store.GetServerPrivateKey(testId); err == nil {
		t.Fatal("Expected error on fetching cleared key")
	}
}

func TestFilesystemKeyStore(t *testing.T) {
	keyDirectory := fmt.Sprintf(".%s%s", string(filepath.Separator), "keys")
	os.MkdirAll(keyDirectory, 0700)
	defer func() {
		os.RemoveAll(keyDirectory)
	}()
	store, err := NewFilesystemKeyStore(keyDirectory)
	if err != nil {
		t.Fatal("error")
	}
	testGeneral(store, t)
	testGeneratingDataEncryptionKeys(store, t)
	testGenerateProxyKeys(store, t)
	testGenerateServerKeys(store, t)
	testReset(store, t)
}
