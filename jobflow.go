/*
 File		:	$File: //depot/stillwater/perfmonitor/jobflow.go $

 Authors	:	E. Theodore L. Omtzigt
 Date		:	5 May 2016

 Source Control Information:
 Version	:	$Revision: #1 $
 Latest		:	$Date: 2016/05/05 $
 Location	:	$Id: //depot/stillwater/perfmonitor/jobflow.go#1 $

 Organization:
		Stillwater Supercomputing, Inc.
		P.O Box 720
		South Freeport, ME 04078-0720

Copyright (c) 2006-2016 E. Theodore L. Omtzigt.  All rights reserved.

Licence      : Stillwater license as defined in this directory

 */
package perfmonitor

import (
	"fmt"
)

type JobFlow struct {
	Arrivals 	uint64
	Busy 		uint64
	Completions	uint64
}

/////////////////////////////////////////////////////////////////
// Selectors

func (f JobFlow) String() string {
	return fmt.Sprintf("[%v,%v,%v]", f.Arrivals, f.Busy, f.Completions)
}

func (f *JobFlow) HasJobFlowBalance() bool {
	return f.Arrivals == f.Completions
}

/////////////////////////////////////////////////////////////////
// Modifiers
