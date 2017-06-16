// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/casbin/casbin"
)

func testGlobalEnforce(t *testing.T, tenant string, sub string, obj string, act string, service string, res bool) {
	var sc SecurityContext
	sc.Tenant = tenant
	sc.Sub = sub
	sc.Obj = obj
	sc.Act = act

	if enforce(sc) != res {
		t.Errorf("%s, %s, %s, %s, %s: %t, supposed to be %t", tenant, sub, obj, act, service, !res, res)
	}
}

func testEnforce(t *testing.T, e *casbin.Enforcer, tenant string, sub string, obj string, act string, service string, res bool) {
	if e.Enforce(tenant, sub, obj, act, service) != res {
		t.Errorf("%s, %s, %s, %s, %s: %t, supposed to be %t", tenant, sub, obj, act, service, !res, res)
	}
}

func TestAdmin(t *testing.T) {
	testGlobalEnforce(t, "admin", "admin", "/", "GET", "nova", true)
	testGlobalEnforce(t, "admin", "admin", "/admin/servers/detail", "GET", "nova", true)
	testGlobalEnforce(t, "admin", "admin", "/admin/extensions", "GET", "nova", true)
	testGlobalEnforce(t, "admin", "admin", "/admin/os-simple-tenant-usage/ce9ff56f5af746de93ec30f387cd7fa8", "GET", "nova", true)
	testGlobalEnforce(t, "admin", "admin", "/admin/flavors/detail", "GET", "nova", true)
	testGlobalEnforce(t, "admin", "admin", "/admin/extensions", "GET", "nova", true)

	testGlobalEnforce(t, "tenant1", "user11", "/admin/servers/detail", "GET", "nova", false)
	testGlobalEnforce(t, "tenant1", "user12", "/admin/servers/detail", "GET", "nova", false)
	testGlobalEnforce(t, "tenant1", "user13", "/admin/servers/detail", "GET", "nova", false)

	testGlobalEnforce(t, "tenant2", "user2", "/admin/servers/detail", "GET", "nova", false)

	testGlobalEnforce(t, "tenant3", "user3", "/admin/servers/detail", "GET", "nova", false)
}

func TestEnable(t *testing.T) {
	e := casbin.NewEnforcer("authz_model.conf", policy_global_enable)

	testEnforce(t, e,"tenant1", "user11", "/metadata", "GET", "patron", true)
	testEnforce(t, e, "tenant1", "user11", "/metadata", "POST", "patron", true)
	testEnforce(t, e, "tenant1", "user11", "/policy", "GET", "patron", true)
	testEnforce(t, e, "tenant1", "user11", "/policy", "POST", "patron", true)

	testEnforce(t, e,"tenant1", "user12", "/metadata", "GET", "patron", false)
	testEnforce(t, e, "tenant1", "user12", "/metadata", "POST", "patron", false)
	testEnforce(t, e, "tenant1", "user12", "/policy", "GET", "patron", false)
	testEnforce(t, e, "tenant1", "user12", "/policy", "POST", "patron", false)

	testEnforce(t, e,"tenant2", "user11", "/metadata", "GET", "patron", false)
	testEnforce(t, e, "tenant2", "user11", "/metadata", "POST", "patron", false)
	testEnforce(t, e, "tenant2", "user11", "/policy", "GET", "patron", false)
	testEnforce(t, e, "tenant2", "user11", "/policy", "POST", "patron", false)
}
