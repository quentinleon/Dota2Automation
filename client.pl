#!/usr/bin/env perl
use strict;
use warnings;

my $port = 1042;
my $ip = "10.10.125.51";
my $dockvol = "/Volumes/AOEU/dota:/dota";
my @args = ("docker", "run", "-d", "-v", $dockvol, "dotatest");
my $containers = 0;
my $peak = 0;
my $finnish = 0;

sub get_response {
	my $done;
	if ($containers < $peak) {
		$done = `echo done | nc $ip $port`;
	} else {
		$done = `echo new | nc $ip $port`;
	}
	return $done;
}

do {
	$containers = int(`docker ps | wc -l`) - 1;
	if (($containers < 4 && !$finnish) || ($containers < $peak)) {
		my $response = get_response();
		if ($response eq "yes") {
			print "Spinning up dota.. ";
			system(@args);
			$containers++;
		} else {
			$finnish = "yes";
		}
	}
	$peak = $containers;
	sleep 2;
} while ($containers > 0);
