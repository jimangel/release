/*
Copyright The Kubernetes Authors.

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

// Code generated by counterfeiter. DO NOT EDIT.
package githubfakes

import (
	"sync"

	githuba "sigs.k8s.io/release-sdk/github"
)

type FakeImpl struct {
	githubStub        func() *githuba.GitHub
	githubMutex       sync.RWMutex
	githubArgsForCall []struct {
	}
	githubReturns struct {
		result1 *githuba.GitHub
	}
	githubReturnsOnCall map[int]struct {
		result1 *githuba.GitHub
	}
	processAssetFilesStub        func([]string) ([]map[string]string, error)
	processAssetFilesMutex       sync.RWMutex
	processAssetFilesArgsForCall []struct {
		arg1 []string
	}
	processAssetFilesReturns struct {
		result1 []map[string]string
		result2 error
	}
	processAssetFilesReturnsOnCall map[int]struct {
		result1 []map[string]string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeImpl) github() *githuba.GitHub {
	fake.githubMutex.Lock()
	ret, specificReturn := fake.githubReturnsOnCall[len(fake.githubArgsForCall)]
	fake.githubArgsForCall = append(fake.githubArgsForCall, struct {
	}{})
	stub := fake.githubStub
	fakeReturns := fake.githubReturns
	fake.recordInvocation("github", []interface{}{})
	fake.githubMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeImpl) GithubCallCount() int {
	fake.githubMutex.RLock()
	defer fake.githubMutex.RUnlock()
	return len(fake.githubArgsForCall)
}

func (fake *FakeImpl) GithubCalls(stub func() *githuba.GitHub) {
	fake.githubMutex.Lock()
	defer fake.githubMutex.Unlock()
	fake.githubStub = stub
}

func (fake *FakeImpl) GithubReturns(result1 *githuba.GitHub) {
	fake.githubMutex.Lock()
	defer fake.githubMutex.Unlock()
	fake.githubStub = nil
	fake.githubReturns = struct {
		result1 *githuba.GitHub
	}{result1}
}

func (fake *FakeImpl) GithubReturnsOnCall(i int, result1 *githuba.GitHub) {
	fake.githubMutex.Lock()
	defer fake.githubMutex.Unlock()
	fake.githubStub = nil
	if fake.githubReturnsOnCall == nil {
		fake.githubReturnsOnCall = make(map[int]struct {
			result1 *githuba.GitHub
		})
	}
	fake.githubReturnsOnCall[i] = struct {
		result1 *githuba.GitHub
	}{result1}
}

func (fake *FakeImpl) processAssetFiles(arg1 []string) ([]map[string]string, error) {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.processAssetFilesMutex.Lock()
	ret, specificReturn := fake.processAssetFilesReturnsOnCall[len(fake.processAssetFilesArgsForCall)]
	fake.processAssetFilesArgsForCall = append(fake.processAssetFilesArgsForCall, struct {
		arg1 []string
	}{arg1Copy})
	stub := fake.processAssetFilesStub
	fakeReturns := fake.processAssetFilesReturns
	fake.recordInvocation("processAssetFiles", []interface{}{arg1Copy})
	fake.processAssetFilesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImpl) ProcessAssetFilesCallCount() int {
	fake.processAssetFilesMutex.RLock()
	defer fake.processAssetFilesMutex.RUnlock()
	return len(fake.processAssetFilesArgsForCall)
}

func (fake *FakeImpl) ProcessAssetFilesCalls(stub func([]string) ([]map[string]string, error)) {
	fake.processAssetFilesMutex.Lock()
	defer fake.processAssetFilesMutex.Unlock()
	fake.processAssetFilesStub = stub
}

func (fake *FakeImpl) ProcessAssetFilesArgsForCall(i int) []string {
	fake.processAssetFilesMutex.RLock()
	defer fake.processAssetFilesMutex.RUnlock()
	argsForCall := fake.processAssetFilesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImpl) ProcessAssetFilesReturns(result1 []map[string]string, result2 error) {
	fake.processAssetFilesMutex.Lock()
	defer fake.processAssetFilesMutex.Unlock()
	fake.processAssetFilesStub = nil
	fake.processAssetFilesReturns = struct {
		result1 []map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) ProcessAssetFilesReturnsOnCall(i int, result1 []map[string]string, result2 error) {
	fake.processAssetFilesMutex.Lock()
	defer fake.processAssetFilesMutex.Unlock()
	fake.processAssetFilesStub = nil
	if fake.processAssetFilesReturnsOnCall == nil {
		fake.processAssetFilesReturnsOnCall = make(map[int]struct {
			result1 []map[string]string
			result2 error
		})
	}
	fake.processAssetFilesReturnsOnCall[i] = struct {
		result1 []map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.githubMutex.RLock()
	defer fake.githubMutex.RUnlock()
	fake.processAssetFilesMutex.RLock()
	defer fake.processAssetFilesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeImpl) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
