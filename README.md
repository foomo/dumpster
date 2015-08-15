# dumpster
dump project data and serve them

## Routes

```
/dumps
    lists dump types
/dumps/<type>
    all dumps of <type>
    Verbs: GET
/dumps/<type>/<id>
    one dump binary download
    VERBS: GET, DELETE, RESTORE
/dumps/<type>/<id>
    VERBS: CREATE
    body: {"id": "foo", "comment": "because we can"}




/remote/<name>/dumps

```
