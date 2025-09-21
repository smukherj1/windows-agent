param (
    # The command to launch.
    [Parameter(Mandatory=$true)]
    [int32]$process_id
)

try {
    # Attempt to get the process by its ID.
    $process = Get-Process `
        -Id $process_id `
        # Throw an exception if the process id doesn't exist.
        -ErrorAction Stop `
    
    # Getting here means the process is currently running.
    Write-Host "true"
}
catch {
    # Getting here means either the process has stopped or
    # the id is invalid.
    Write-Host "false"
}