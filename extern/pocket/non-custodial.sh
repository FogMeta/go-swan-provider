#!/usr/bin/expect

if {$argc < 9} {
  puts "Usage: non-custodial.sh <pubKey> <outputAddr> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore> <passwd>"
  exit 1
}

set pubKey [lindex $argv 0]
set outputAddr [lindex $argv 1]
set amount [lindex $argv 2]
set relayChainIDs [lindex $argv 3]
set serviceURI [lindex $argv 4]
set networkID [lindex $argv 5]
set fee [lindex $argv 6]
set isBefore [lindex $argv 7]
set passwd [lindex $argv 8]

set command "pocket nodes stake custodial $pubKey $outputAddr $amount $relayChainIDs $serviceURI $networkID $fee $isBefore"
spawn sh -c "echo $command"

$command
sleep 1
send -- "$passwd\n"

expect eof
exit