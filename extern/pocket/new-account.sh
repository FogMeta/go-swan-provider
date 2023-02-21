#!/usr/bin/expect

if {$argc < 2} {
  puts "Usage: new-account.sh <passwd>"
  exit 1
}

set passwd [lindex $argv 0]

set command "pocket accounts create"
spawn sh -c "echo $command"

$command
sleep 1
send -- "$passwd\n"
sleep 1
send -- "$passwd\n"

expect eof
exit