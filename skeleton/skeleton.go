/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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

package skeleton

import (
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	Name = "skeleton"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

/*
	The interface for collectPlugins defines 3 functions:

	GetConfigPolicy() (*cpolicy.ConfigPolicy, error)
	CollectMetrics([]MetricType) ([]MetricType, error)
	GetMetricTypes(PluginConfigType) ([]MetricType, error)

	This check will error at compile time if the interface is not implemented.
*/
var _ plugin.CollectorPlugin = (*Skeleton)(nil)

// Our type that implements plugin.CollectorPlugin interface
// defined in intel-sdi-x/snap/control/plugin
type Skeleton struct {
}

/*
	CollectMetrics() will be called by snap when a task that collects one of the metrics returned from this plugins
	GetMetricTypes() is started. The input will include a slice of all the metric types being collected.

	The output is the collected metrics as plugin.MetricType's and an error.
*/
func (f *Skeleton) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	// Create a slice to hold our results
	results := []plugin.MetricType{}
	// Go through all metrics we were asked to collect
	for _, metric := range mts {
		// call the helper function to collect an individal metric
		// Returns a slice since a dynamic metric could return more than one result.
		// For example: collecting /foo/* would return all metrics that start with '/foo'
		metrics, err := f.collectMetric(metric)
		// handle the error as appropriate
		if err != nil {
			// Do error handling
		}
		results = append(results, metrics...)
	}
	return results, nil
}

/*
	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (f *Skeleton) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	mts = append(mts, plugin.MetricType{Namespace_: core.NewNamespace("foo", "bar", "baz")})
	return mts, nil
}

/*
	GetConfigPolicy() returns the configPolicy for your plugin.

	A config policy is how users can provide configuration info to
	plugin. Here you define what sorts of config info your plugin
	needs and/or requires.
*/
func (f *Skeleton) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	/*
		There are 4 different rule types:

		StringRule
		IntegerRule
		BoolRule
		FloatRule

		They are all defined following a similar pattern:

		rule, err := cpolicy.NewStringRule(<name of field>, <required (true/false)>, <default value>)

		rules need to be added to Policy Nodes. A policy node can be added to the config policy at either the root level,
		or specific to a certain metric.
	*/
	// Add a string rule for name, it is not required and defaults to bob
	rule, err := cpolicy.NewStringRule("name", false, "bob")
	if err != nil {
		return nil, err
	}
	// Add a stringrule for password, it is required and has no default
	rule2, err := cpolicy.NewStringRule("password", true)
	if err != nil {
		return nil, err
	}
	// Create a cpolicy.PolicyNode and add our rules to it
	p := cpolicy.NewPolicyNode()
	p.Add(rule)
	p.Add(rule2)
	// Add this policy node to apply only to namespaces that begin with /foo/bar/baz
	c.Add([]string{"foo", "bar", "baz"}, p)

	return c, nil
}

//Meta returns meta data for use in main.go
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		// Defined at the top of the file
		Name,
		Version,
		Type,
		// Content types that this plugin supports. Not recommended to use other
		// content types unless there are specific constraints that require it.
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		// Is this plugin signed?
		plugin.Unsecure(true),
		// Routing strategy options are currently:
		//
		// DefaulRouting - A collectMetrics request may be initiated on any one of my running
		// plugins (snap may start more than one instance of your plugin based on demand).
		//
		// StickyRouting - A task is tied to a specific plugin (useful if you have state that cannot
		// be shared between different tasks like a database connection or other configurable option)
		plugin.RoutingStrategy(plugin.DefaultRouting),
		// The amount of time that a metric is valid inside the cache for. If your plugin is called for the same
		// metric in an interval < the CacheTTL it will use cached values instead of calling the plugin again.
		plugin.CacheTTL(1100*time.Millisecond),
	)
}

// This is not apart of the interface required to imeplement plugin.CollectorPlugin. It is just a helper function to
// process one metric at a time and return the result(s).
//
// Because a dynamic metric (defined with a * in the namespace like /bin/foo/*/stat) can return more
// than one result we return a slice here.
func (f *Skeleton) collectMetric(metric plugin.MetricType) ([]plugin.MetricType, error) {
	//For each metric we need to decide what data we want to return. We do this by examining the namespace
	// and figuring out which of the MetricTypes we defined in GetMetricTypes this metric metric is targeting.
	//
	// This is where the bulk of the work for a collect plugin is done.
	//
	// For example purposese we will return a static value here.
	mt := plugin.MetricType{
		Namespace_: metric.Namespace(),
		Version_:   metric.Version(),
		// Data_ is type interface{} so any type can be assigned to it.
		Data_:      1,
		Timestamp_: time.Now(),
	}

	return []plugin.MetricType{mt}, nil
}
