# farchive

Manage hash in sqlite3, as an personal alternative to [snapraid content](https://www.snapraid.it/).

run

    λ farchive run
    2024/06/06 18:17:31 new file go.sum
    2024/06/06 18:17:31 new file go.mod
    2024/06/06 18:17:31 new file main.go

    λ sqlite3 farchive.db 'select * from file'
    -- Loading resources from /home/rok/.sqliterc
    path        abs                                               size  hash              modifiedAt  validatedAt
    ----------  ------------------------------------------------  ----  ----------------  ----------  -----------
    go.sum      /home/rok/src/github.com/aca/farchive/go.sum      6486  340ea90ffa909eea  1717664739  1717665451
    go.mod      /home/rok/src/github.com/aca/farchive/go.mod      1078  b71ed6c8a8a82c7c  1717664739  1717665451
    main.go     /home/rok/src/github.com/aca/farchive/main.go     3185  48ed034ab7d791a1  1717665448  1717665451

    # Edit main.go
    λ farchive run
    2024/06/06 18:18:26 UPDATE HASH main.go 48ed034ab7d791a1 64291ff9f1b0960d


diff

    # after disk sync, check if there is any difference
    λ farchive diff /mnt/disk0/farchive.db /mnt/disk1/farchive.db

