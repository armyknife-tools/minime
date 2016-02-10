#!/usr/bin/env perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Generate system call table for Darwin from sys/syscall.h

use strict;

my $command = "mksysnum_darwin.pl " . join(' ', @ARGV);

print <<EOF;
// $command
// MACHINE GENERATED BY THE ABOVE COMMAND; DO NOT EDIT

package unix

const (
EOF

while(<>){
	if(/^#define\s+SYS_(\w+)\s+([0-9]+)/){
		my $name = $1;
		my $num = $2;
		$name =~ y/a-z/A-Z/;
		print "	SYS_$name = $num;"
	}
}

print <<EOF;
)
EOF
