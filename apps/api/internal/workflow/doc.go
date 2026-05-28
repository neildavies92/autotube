// Package workflow defines worker-ready boundaries for background work.
//
// The durable queue is intentionally not implemented here. These contracts let
// API handlers, future worker processes, and persistence code agree on job
// shapes before a concrete queue such as River is introduced.
package workflow
