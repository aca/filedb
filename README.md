# farchive

Manage hash in sqlite3, as an alternative to "snapraid content".

    λ go run .
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

    λ go run .
    2024/06/06 18:18:26 UPDATE HASH main.go 48ed034ab7d791a1 64291ff9f1b0960d
