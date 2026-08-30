#!/usr/bin/env perl

use strict;
use warnings;
use Getopt::Long;

# Script to collect and export Android hardware telemetry metrics
my $serial = '';
my $format = 'table'; # table or json
my $help   = 0;

GetOptions(
    'serial=s' => \$serial,
    'format=s' => \$format,
    'help'     => \$help
);

if ($help) {
    print "Usage: $0 [--serial <serial>] [--format <table|json>]\n";
    print "Collect and format Android device hardware and system telemetry\n";
    exit 0;
}

# Run command helper
sub run_adb {
    my ($cmd) = @_;
    my $full_cmd = "adb";
    if ($serial ne '') {
        $full_cmd .= " -s " . quotemeta($serial);
    }
    $full_cmd .= " $cmd";
    my $out = `$full_cmd 2>/dev/null`;
    return $out || '';
}

# 1. Device identification
my $model   = run_adb("shell getprop ro.product.model");
my $product = run_adb("shell getprop ro.product.name");
my $os_ver  = run_adb("shell getprop ro.build.version.release");
my $sdk_ver = run_adb("shell getprop ro.build.version.sdk");
chomp($model, $product, $os_ver, $sdk_ver);

# 2. Battery info
my $battery_out = run_adb("shell dumpsys battery");
my $battery_level = "Unknown";
my $battery_temp  = "Unknown";
my $battery_volt  = "Unknown";
my $battery_stat  = "Unknown";

if ($battery_out =~ /level:\s*(\d+)/) {
    $battery_level = "$1%";
}
if ($battery_out =~ /temperature:\s*(\d+)/) {
    $battery_temp = sprintf("%.1f C", $1 / 10.0);
}
if ($battery_out =~ /voltage:\s*(\d+)/) {
    $battery_volt = sprintf("%.2f V", $1 / 1000.0);
}
if ($battery_out =~ /status:\s*(\d+)/) {
    my %status_map = (2 => "Charging", 3 => "Discharging", 4 => "Not charging", 5 => "Full");
    $battery_stat = $status_map{$1} || "Unknown";
}

# 3. Memory info
my $mem_out = run_adb("shell cat /proc/meminfo");
my $mem_total = "Unknown";
my $mem_avail = "Unknown";

if ($mem_out =~ /MemTotal:\s*(\d+)\s*kB/) {
    $mem_total = sprintf("%.1f GB", $1 / (1024 * 1024));
}
if ($mem_out =~ /MemAvailable:\s*(\d+)\s*kB/) {
    $mem_avail = sprintf("%.1f GB", $1 / (1024 * 1024));
}

# 4. Storage info
my $df_out = run_adb("shell df -h /data");
my $storage_size = "Unknown";
my $storage_free = "Unknown";
my $storage_used = "Unknown";

for my $line (split(/\n/, $df_out)) {
    next if $line =~ /Filesystem/ || $line =~ /^\s*$/;
    my @parts = split(/\s+/, $line);
    if (@parts >= 5) {
        $storage_size = $parts[1];
        $storage_used = $parts[2];
        $storage_free = $parts[3];
        last;
    }
}

if (lc($format) eq 'json') {
    printf("{\n");
    printf(qq{  "model": "%s",\n}, $model);
    printf(qq{  "product": "%s",\n}, $product);
    printf(qq{  "android_version": "%s",\n}, $os_ver);
    printf(qq{  "sdk_level": "%s",\n}, $sdk_ver);
    printf(qq{  "battery": {\n    "level": "%s",\n    "status": "%s",\n    "temperature": "%s",\n    "voltage": "%s"\n  },\n},
        $battery_level, $battery_stat, $battery_temp, $battery_volt);
    printf(qq{  "memory": {\n    "total": "%s",\n    "available": "%s"\n  },\n}, $mem_total, $mem_avail);
    printf(qq{  "storage": {\n    "size": "%s",\n    "used": "%s",\n    "free": "%s"\n  }\n},
        $storage_size, $storage_used, $storage_free);
    printf("}\n");
} else {
    print "========================================\n";
    print "   Android Device Telemetry Metrics     \n";
    print "========================================\n";
    printf("%-20s: %s (%s)\n", "Model / Product", $model, $product);
    printf("%-20s: Android %s (API %s)\n", "OS Version", $os_ver, $sdk_ver);
    printf("%-20s: %s (%s, %s, %s)\n", "Battery", $battery_level, $battery_stat, $battery_temp, $battery_volt);
    printf("%-20s: %s avail / %s total\n", "RAM Memory", $mem_avail, $mem_total);
    printf("%-20s: %s free / %s total (used: %s)\n", "Internal Storage", $storage_free, $storage_size, $storage_used);
    print "========================================\n";
}

