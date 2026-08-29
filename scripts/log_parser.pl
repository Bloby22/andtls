#!/usr/bin/env perl

use strict;
use warnings;
use Getopt::Long;

# ANSI color definitions for terminal formatting
my %COLORS = (
    'V'     => "\e[0;37m",  # Verbose: White
    'D'     => "\e[0;34m",  # Debug: Blue
    'I'     => "\e[0;32m",  # Info: Green
    'W'     => "\e[1;33m",  # Warning: Yellow
    'E'     => "\e[1;31m",  # Error: Red
    'F'     => "\e[1;35m",  # Fatal: Magenta
    'TAG'   => "\e[1;36m",  # Tag: Cyan
    'PID'   => "\e[0;35m",  # PID: Purple
    'RESET' => "\e[0m"
);

my $tag_filter = '';
my $min_level  = 'V';
my $serial     = '';
my $help       = 0;

GetOptions(
    'tag=s'    => \$tag_filter,
    'level=s'  => \$min_level,
    'serial=s' => \$serial,
    'help'     => \$help
);

if ($help) {
    print "Usage: $0 [--serial <serial>] [--tag <regex>] [--level <V|D|I|W|E|F>]\n";
    print "Colorized stream parser for Android logcat output\n";
    exit 0;
}

my %LEVEL_PRIORITY = (
    'V' => 1,
    'D' => 2,
    'I' => 3,
    'W' => 4,
    'E' => 5,
    'F' => 6
);

my $min_priority = $LEVEL_PRIORITY{uc($min_level)} || 1;

# Build adb command
my $cmd = "adb";
if ($serial ne '') {
    $cmd .= " -s " . quotemeta($serial);
}
$cmd .= " logcat -v time";

print "\e[1;34mStarting colorized logcat parser (Level: >=$min_level, Tag: '$tag_filter')...\e[0m\n";

open(my $pipe, "-|", $cmd) or die "Failed to execute adb logcat: $!\n";

while (my $line = <$pipe>) {
    chomp($line);

    # Standard logcat time format: "08-22 19:30:15.123 D/TagName( 1234): Message text"
    if ($line =~ /^(\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\.\d+)\s+([VDIWEF])\/(.*?)\(\s*(\d+)\):\s*(.*)$/) {
        my ($time, $level, $tag, $pid, $msg) = ($1, $2, $3, $4, $5);

        my $priority = $LEVEL_PRIORITY{$level} || 1;
        next if $priority < $min_priority;

        if ($tag_filter ne '' && $tag !~ /$tag_filter/i) {
            next;
        }

        my $lvl_color = $COLORS{$level} || $COLORS{'RESET'};
        my $tag_color = $COLORS{'TAG'};
        my $pid_color = $COLORS{'PID'};
        my $reset     = $COLORS{'RESET'};

        printf("%s %s[%s]%s %s%-20s%s (%s%5s%s): %s%s%s\n",
            $time,
            $lvl_color, $level, $reset,
            $tag_color, substr($tag, 0, 20), $reset,
            $pid_color, $pid, $reset,
            $lvl_color, $msg, $reset
        );
    } else {
        # Raw fallback line
        print "$line\n";
    }
}

close($pipe);

