default:
    @echo 'Usage: just run-[os]'
    @echo 'Example: '
    @just -l

run-bash:
    while ! go run . ; do :; done

#run-powershell:
