#!/usr/bin/expect

if {$argc < 8} {
  puts "Usage: custodial.sh <address> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore> <passwd>"
  exit 1
}

set address [lindex $argv 0]
set amount [lindex $argv 1]
set relayChainIDs [lindex $argv 2]
set serviceURI [lindex $argv 3]
set networkID [lindex $argv 4]
set fee [lindex $argv 5]
set isBefore [lindex $argv 6]
set passwd [lindex $argv 7]

set command "pocket nodes stake custodial $address $amount $relayChainIDs $serviceURI $networkID $fee $isBefore"
spawn sh -c "$command"
sleep 1
send -- "$passwd\n"

expect eof
exit