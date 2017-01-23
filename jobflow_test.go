/*
 File		:	$File: //depot/stillwater-sc/perfmonitor/jobflow_test.go $

 Authors	:	E. Theodore L. Omtzigt
 Date		:	5 May 2016

 Source Control Information:
 Version	:	$Revision: #1 $
 Latest		:	$Date: 2016/05/05 $
 Location	:	$Id: //depot/stillwater-sc/perfmonitor/jobflow_test.go#1 $

 Organization:
		Stillwater Supercomputing, Inc.
		P.O Box 720
		South Freeport, ME 04078-0720

Copyright (c) 2006-2016 E. Theodore L. Omtzigt.  All rights reserved.

Licence      : Stillwater license as defined in this directory
 */
package perfmonitor

import (
	"testing"
)

/////////////////////////////////////////////////////////////////
// Test cases

func TestJobFlow_JobFlowBalance(t *testing.T) {
	var jf1 JobFlow = JobFlow{11,0,10}  // no job flow balance
	var jf2 JobFlow = JobFlow{10,0,10}  // does have balance
	if jf1.HasJobFlowBalance() || !jf2.HasJobFlowBalance()  {
		t.Errorf("JobFlow balance test is incorrect")
	}
}

/////////////////////////////////////////////////////////////////
// Examples


/////////////////////////////////////////////////////////////////
// Benchmarks
