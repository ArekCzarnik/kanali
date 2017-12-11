// Copyright (c) 2017 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package steps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/plugin"
	opentracing "github.com/opentracing/opentracing-go"
	"k8s.io/client-go/informers/core"
)

type pluginsOnRequestStep struct{}

func NewPluginsOnRequestStep() Step {
	return pluginsOnRequestStep{}
}

func (step pluginsOnRequestStep) GetName() string {
	return "plugin onrequest"
}

func (step pluginsOnRequestStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := getPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnRequest(ctx, m, plugin, *proxy, r, trace, *p); err != nil {
			return err
		}
	}
	return nil
}

func doOnRequest(ctx context.Context, m *metrics.Metrics, plugin v2.Plugin, proxy v2.ApiProxy, req *http.Request, span opentracing.Span, p plugin.Plugin) (e error) {
	logger := logging.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("OnRequest paniced: %v", r))
			e = errors.New("OnRequest paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_REQUEST: %s", plugin.Name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnRequest(ctx, plugin.Config, m, proxy, req, sp)
}