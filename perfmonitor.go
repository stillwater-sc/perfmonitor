/**
 * File		:	$File: //depot/stillwater/perfmonitor/perfmonitor.go $
 *
 * Authors	:	E. Theodore L. Omtzigt
 * Date		:	5 May 2016
 *
 * Source Control Information:
 * Version	:	$Revision: #1 $
 * Latest	:	$Date: 2016/05/05 $
 * Location	:	$Id: //depot/stillwater/perfmonitor/perfmonitor.go#1 $
 *
 * Organization:
 *		Stillwater Supercomputing, Inc.
 *		P.O Box 720
 *		South Freeport, ME 04078-0720
 *
 * Copyright (c) 2006-2016 E. Theodore L. Omtzigt.  All rights reserved.
 *
 * Licence      : Stillwater license as defined in this directory
 *
 * The PerfMonitor is an infrastructure that tracks Operational Analysis attributes
 * of resources, and time series sequences of transactions. The goal of the PerfMonitor
 * is to deliver performance metrics for quality assurance and regression testing.
  *
 * The transaction monitoring should be capable to trace the update sequences of a
 * system of recurrence equations. Domain flow algorithms can be viewed as injecting data
 * elements into domains of computation that evolve intermediate values till they are
 * ejected back into memory, or another domain. For understanding the dynamics and
 * functional validity of a domain flow algorithm, tracing these evolutions is valuable.
 *
 * The general structure of a Domain Flow Algorithm is:
 * input ((i,j,k) | l <= i,j <= u, k = c) {
 *		// injection of a memory data structure into a Domain of Computation
 * 	a(i,j,k) = A(i,j)
 *	b(k,j,i) = B(i,j)
 * }
 * compute ((i,j,k) | l <= i,j,k <= u) {
 *		// evolution of a computation inside a Domain of Computation
 * 	a(i,j,k) = f(a(i,j,k-1))
 *	b(i,j,k) = g(b(i-1,j,k))
 * }
 * output ((i,j,k) | l <= i,j <= u, k = d) {
 *		// ejection of a memory data structure from a Domain of Computation
 *	Aprime(i,j) = a(i,j,k)
 * }
 *
 * The transaction tracing would tag an input data element from the input() domain,
 * for example, A(1,1).
 * Then trace the evolution of the recurrence; a(i,j,k) = f(a(i,j,k-1)) through the domain.
 *
 * The transaction traces generated would look like this:
 * A(1,1) (t0,R_a,v1) (t1,R_b,v2) (t2,R_c,v3) (t3,R_d,v4) ... etc.
 * A(1,2) (t0,R_w,v1) (t1,R_x,v2) (t2,R_y,v3) (t3,R_z,v4) ... etc.
 * A(2,1) (t0,R_a,v1) (t1,R_b,v2) (t2,R_c,v3) (t3,R_d,v4) ... etc.
 * t0, t1, t2, etc. are time stamps in nanoseconds.
 * R_a, R_b, R_c, etc. are resource identifiers
 * v1, v2, v3, etc. are the values of the intermediate values of the recurrence equation.
 *
 * We would like to be able to leverage the InfluxDB time series database for data
 * storage, management, and queries, so we are using the same line protocol as that
 * database. This would allow us to replay the events captured in the PerfMonitor
 * and send them to an InfluxDB instance. We have an in-memory PerfMonitor in the
 * simulator because of the performance requirements. InfluxDB replays would only be
 * required for inspection and debug. Most of the time, the low level operational
 * analysis attributes will be enough to validate results, and support regression
 * testing.
 *
 * The InfluxDB line protocol is structured as follows:
 *
{
    "database": "foo",
    "retentionPolicy": "bar",
    "points": [
        {
            "name": "measurement",
            "tags": {
                "host": "server01",
                "region": "us-east1",
                "tag1": "value1",
                "tag2": "value2",
                "tag2": "value3",
                "tag4": "value4",
                "tag5": "value5",
                "tag6": "value6",
                "tag7": "value7",
                "tag8": "value8"
            },
            "time": 14244733039069373,
            "precision": "n",
            "fields": {
                    "value": 4541770385657154000
            }
        }
    ]
}

The measurement name, a set of tags, a timestamp plus precision identifier, and a
set of fields of values defines a point in the time series. The InfluxDB database
indexes on measurement and the tags, but not on the fields.

We could map our computational events on this model like so:
A point:
 {
            "name": "A(1,1)",
            "tags": {
                "pe": "[1][1]",			// the resource tag
                "lp": "(1,1,1)",		// lattice point of the computational event
                "re": "a"			// recurrence id
            },
            "time": 1,				// clock ticks in terms of nsec
            "precision": "n",
            "fields": {
                    "value": 1.0f
            }
  },
 {
            "name": "A(1,1)",
            "tags": {
                "pe": "[1][2]",			// the resource tag
                "lp": "(1,1,2)",		// lattice point of the computational event
                "re": "a"			// recurrence id
            },
            "time": 1,				// clock ticks in terms of nsec
            "precision": "n",
            "fields": {
                    "value": 2.0f
            }
  },
 {
            "name": "A(1,1)",
            "tags": {
                "pe": "[1][3]",			// the resource tag
                "lp": "(1,1,3)",		// lattice point of the computational event
                "re": "a"			// recurrence id
            },
            "time": 1,				// clock ticks in terms of nsec
            "precision": "n",
            "fields": {
                    "value": 3.0f
            }
  },
 */
package perfmonitor

import (
	"fmt"
	"github.com/golang/glog"
)

type ResourceTag uint64
type ResourceAddress [3]int
type PerfMonitor map[ResourceTag]*JobFlow

/////////////////////////////////////////////////////////////////
// Selectors

func (t *ResourceTag) GenerateUniqueIdentifier(index ResourceAddress) ResourceTag {
	*t = ResourceTag(uint64(index[0]) | uint64(index[1]) << 16 | uint64(index[2]) << 32)
	return *t
}

func (oa *PerfMonitor) Print() {
	for key, value := range (*oa) {
		fmt.Printf("id[%#x] = %v\n", key, value)
	}
}

/////////////////////////////////////////////////////////////////
// Modifiers
func (oa *PerfMonitor) Arrival(tag ResourceTag, timeStamp uint64) {
	// TODO: protect from concurrent go-routines
	var jf *JobFlow
	jf, ok := (*oa)[tag]
	if ok {
		jf.Arrivals = jf.Arrivals + 1
	} else {
		jf = new(JobFlow)
		jf.Arrivals = 1
		(*oa)[tag] = jf
	}
}

func (oa *PerfMonitor) Completion(tag ResourceTag, timeStamp uint64) {
	// TODO: protect from concurrent go-routines
	var jf *JobFlow
	jf, ok := (*oa)[tag]
	if ok {
		jf.Completions = jf.Completions + 1
	} else {
		// this is bad, Arrival should have allocated a new JobFlow
		glog.Errorf("Completion for tag %d was called before Arrival", tag)

	}
}

