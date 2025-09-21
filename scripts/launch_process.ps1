# Launch a process asynchronously

param (
    # The command to launch.
    [Parameter(Mandatory=$true)]
    [string]$cmd,
    
    # Arguments to pass to the command.
    [string]$cmd_arguments=""
)

# Start with the mandatory arguments
$runnerArgs = @(
    "-file"
    "run_process.ps1"
    "-cmd"
    $cmd
)

# Conditionally add the optional argument if it has a value
if ($cmd_arguments) {
    $runnerArgs += "-cmd_arguments"
    $runnerArgs += $cmd_arguments
}

$processArgs = @{
    FilePath = "powershell"
    PassThru = $true
    Wait = $false
    NoNewWindow = $false
    RedirectStandardError = "C:/Apps/launch_process.stderr.txt"
    RedirectStandardOutput = "C:/Apps/launch_process.stdout.txt"
    ErrorAction = "Stop"
    ArgumentList = $runnerArgs
}

$process = Start-Process @processArgs

$agent_pid = $process.Id

Write-Host "$agent_pid"