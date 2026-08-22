#!/usr/bin/env perl

use strict;
use warnings;

# Perl script to generate Android udev rules for Linux systems
my %KNOWN_VENDORS = (
    '18d1' => 'Google / Pixel',
    '04e8' => 'Samsung',
    '2717' => 'Xiaomi',
    '0b05' => 'ASUS',
    '0fce' => 'Sony',
    '22b8' => 'Motorola',
    '12d1' => 'Huawei',
    '2a70' => 'OnePlus',
    '17ef' => 'Lenovo',
    '0bb4' => 'HTC',
    '1004' => 'LG',
    '2ae5' => 'Fairphone',
    '2e04' => 'Nothing'
);

print "========================================\n";
print " Android udev Rules Generator for Linux \n";
print "========================================\n\n";

# Check current user in plugdev or adbusers group
my $user = $ENV{'USER'} || 'user';
my $rules_path = "/etc/udev/rules.d/51-android.rules";

print "Checking connected USB devices with lsusb...\n";
my @detected_vendors;

if (open(my $fh, "-|", "lsusb 2>/dev/null")) {
    while (my $line = <$fh>) {
        if ($line =~ /ID\s+([0-9a-fA-F]{4}):[0-9a-fA-F]{4}\s+(.*)/) {
            my $vid = lc($1);
            my $desc = $2;
            if (exists $KNOWN_VENDORS{$vid}) {
                push @detected_vendors, { id => $vid, name => $KNOWN_VENDORS{$vid}, desc => $desc };
                print "Found known Android device: $KNOWN_VENDORS{$vid} (ID: $vid) - $desc\n";
            }
        }
    }
    close($fh);
} else {
    print "Note: lsusb command not available, generating complete rule set\n";
}

print "\nGenerated udev rules content for $rules_path:\n";
print "--------------------------------------------------------\n";

my $header = << "EOF";
# 51-android.rules - Android ADB USB device permissions
# Place this file in /etc/udev/rules.d/51-android.rules
# Reload with: sudo udevadm control --reload-rules && sudo udevadm trigger

EOF

print $header;

for my $vid (sort keys %KNOWN_VENDORS) {
    my $name = $KNOWN_VENDORS{$vid};
    printf(qq{SUBSYSTEM=="usb", ATTR{idVendor}=="%s", MODE="0660", GROUP="adbusers", TAG+="uaccess" # %s\n}, $vid, $name);
}

print "--------------------------------------------------------\n";
print "\nInstallation steps:\n";
print "1. Save the above content to $rules_path\n";
print "   sudo perl $0 | grep 'SUBSYSTEM' | sudo tee $rules_path\n";
print "2. Add your user to adbusers or plugdev group:\n";
print "   sudo usermod -aG adbusers $user\n";
print "3. Reload udev daemon:\n";
print "   sudo udevadm control --reload-rules && sudo udevadm trigger\n";

