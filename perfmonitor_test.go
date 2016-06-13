/**
 * File		:	$File: //depot/stillwater/perfmonitor/perfmonitor_test.go $
 *
 * Authors	:	E. Theodore L. Omtzigt
 * Date		:	5 May 2016
 *
 * Source Control Information:
 * Version	:	$Revision: #1 $
 * Latest	:	$Date: 2016/05/05 $
 * Location	:	$Id: //depot/stillwater/perfmonitor/perfmonitor_test.go#1 $
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
 */
package perfmonitor

import (
	"testing"
	"fmt"
	"strings"
)

/////////////////////////////////////////////////////////////////
// Test cases

func TestOperationalQuantities_String(t *testing.T) {
	var i uint64
	var perfmon PerfMonitor = make(map[ResourceTag]*JobFlow)
	for i = 0; i < 11; i++ {
		perfmon.Arrival(0, i)
	}
	for i = 0; i < 10; i++ {
		perfmon.Completion(0, i+5)
	}
	var jf *JobFlow = perfmon[0]
	a := fmt.Sprintf("%s",jf.String())
	if strings.Compare(a, "[11,0,10]") != 0 {
		t.Errorf("String conversion is incorrect")
	}
}

func TestPerfMonitor_OperationalAnalysis(t *testing.T) {

}
/////////////////////////////////////////////////////////////////
// Examples


/////////////////////////////////////////////////////////////////
// Benchmarks
