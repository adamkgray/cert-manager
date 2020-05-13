/*
Copyright 2020 The Jetstack cert-manager contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package accounts

import (
	"testing"

	cmacme "github.com/jetstack/cert-manager/pkg/apis/acme/v1alpha2"
	"github.com/jetstack/cert-manager/pkg/util/pki"
)

func TestRegistry_AddClient(t *testing.T) {
	r := NewDefaultRegistry()
	pk, err := pki.GenerateRSAPrivateKey(2048)
	if err != nil {
		t.Fatal(err)
	}

	// Register a new client
	r.AddClient("abc", cmacme.ACMEIssuer{}, pk)

	c, err := r.GetClient("abc")
	if err != nil {
		t.Errorf("unexpected error getting client: %v", err)
	}
	if c == nil {
		t.Error("nil client returned")
	}
}

func TestRegistry_RemoveClient(t *testing.T) {
	r := NewDefaultRegistry()
	pk, err := pki.GenerateRSAPrivateKey(2048)
	if err != nil {
		t.Fatal(err)
	}

	// Register a new client
	r.AddClient("abc", cmacme.ACMEIssuer{}, pk)

	c, err := r.GetClient("abc")
	if err != nil {
		t.Errorf("unexpected error getting client: %v", err)
	}
	if c == nil {
		t.Error("nil client returned")
	}

	r.RemoveClient("abc")
	c, err = r.GetClient("abc")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound but got: %v", err)
	}
	if c != nil {
		t.Error("expected nil client to be returned")
	}
}

func TestRegistry_RemoveClient_EmptyRegistry(t *testing.T) {
	r := NewDefaultRegistry()
	r.RemoveClient("abc")
	c, err := r.GetClient("abc")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound but got: %v", err)
	}
	if c != nil {
		t.Error("expected nil client to be returned")
	}
}

func TestRegistry_ListClients(t *testing.T) {
	r := NewDefaultRegistry()
	pk, err := pki.GenerateRSAPrivateKey(2048)
	if err != nil {
		t.Fatal(err)
	}

	// Register a new client
	r.AddClient("abc", cmacme.ACMEIssuer{}, pk)
	l := r.ListClients()
	if len(l) != 1 {
		t.Errorf("expected ListClients to have 1 item but it has %d", len(l))
	}

	// Register a second client
	r.AddClient("abc2", cmacme.ACMEIssuer{}, pk)
	l = r.ListClients()
	if len(l) != 2 {
		t.Errorf("expected ListClients to have 2 items but it has %d", len(l))
	}

	// Register a third client with the same options as the second, meaning
	// it should be de-duplicated
	r.AddClient("abc2", cmacme.ACMEIssuer{}, pk)
	l = r.ListClients()
	if len(l) != 2 {
		t.Errorf("expected ListClients to have 2 items but it has %d", len(l))
	}

	// Update the second client with a new server URL
	r.AddClient("abc2", cmacme.ACMEIssuer{Server: "abc.com"}, pk)
	l = r.ListClients()
	if len(l) != 2 {
		t.Errorf("expected ListClients to have 2 items but it has %d", len(l))
	}
}

func TestRegistry_AddClient_UpdatesExistingWhenPrivateKeyChanges(t *testing.T) {
	r := NewDefaultRegistry()
	pk, err := pki.GenerateRSAPrivateKey(2048)
	if err != nil {
		t.Fatal(err)
	}
	pk2, err := pki.GenerateRSAPrivateKey(2048)
	if err != nil {
		t.Fatal(err)
	}

	// Register a new client
	r.AddClient("abc", cmacme.ACMEIssuer{}, pk)
	l := r.ListClients()
	if len(l) != 1 {
		t.Errorf("expected ListClients to have 1 item but it has %d", len(l))
	}

	// Update the client with a new private key
	r.AddClient("abc", cmacme.ACMEIssuer{}, pk2)
	l = r.ListClients()
	if len(l) != 1 {
		t.Errorf("expected ListClients to have 1 item but it has %d", len(l))
	}
}
