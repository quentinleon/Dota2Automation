#!/usr/bin/env perl
use strict;
use warnings;

if ($#ARGV != 2) {
	die "ERR: usage: ./client.pl ip port docker:volume\n";
}

# we sleep 5 seconds becuase quentin slow
sleep 5;

my $ip = $ARGV[0];
my $port = int($ARGV[1]);
my $dockvol = $ARGV[2];
my @args = ("/usr/local/bin/docker", "run", "-d", "-v", $dockvol, "dotatest");
my $containers = 0;
my $peak = 0;
my $finnish = 0;

sub get_response {
	my $done;
	if ($containers < $peak) {
		$done = `echo done | nc $ip $port`
			or die "no response from server\n";
	} else {
		$done = `echo new | nc $ip $port`
			or die "no response from server\n";
	}
	return $done;
}

do {
	$containers = int(`/usr/local/bin/docker ps | wc -l`) - 1;
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
