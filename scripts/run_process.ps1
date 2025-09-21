# Lower-level synchronous process runner. Use launch_process.ps1 which wraps
# around this script to launch a process asynchronously.

# Launches a process in the same window and waits for it to complete and then
# prints the exit code to stdout.
param (
    # The command to launch.
    [Parameter(Mandatory=$true)]
    [string]$cmd,
    
    # Arguments to pass to the command.
    [string]$cmd_arguments=""
)

$processArgs = @{
    FilePath = $cmd
    PassThru = $true
    Wait = $true
    NoNewWindow = $true
    RedirectStandardError = "C:/Apps/stderr.txt"
    RedirectStandardOutput = "C:/Apps/stdout.txt"
    ErrorAction = "Stop"
}

# ArgumentList is the list of arguments that'll be passed to the launched
# command. ArgumentList requires the given list of arguments to be
# non-empty.
if ($cmd_arguments) {
    $processArgs.Add("ArgumentList", $cmd_arguments)
}

$process = Start-Process @processArgs


$agent_pid = $process.Id

Write-Host "pid: $agent_pid"
Write-Host "exit_code: $($process.ExitCode)"
Write-Host "exit_time: $($process.ExitTime)"