#!/usr/bin/expect

if {$argc < 2} {
  puts "Usage: set-validator.sh <address> <passwd>"
  exit 1
}

set address [lindex $argv 0]
set passwd [lindex $argv 1]

set command "pocket accounts set-validator $address"
spawn sh -c "$command"
sleep 1
send -- "$passwd\n"

expect eof
exit