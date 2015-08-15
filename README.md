# dumpster
dump project data and serve them

## Routes

```
/dumps
    lists dump types
/dumps/<type>
    all dumps of <name>
    Verbs: GET
/dumps/<type>/<id>
    one dump
    VERBS: GET, DELETE, RESTORE
/dumps/<type>/<id>
    VERBS: CREATE
    body: {"id": "foo", "comment": "because we can"}




/remote/<name>/dumps

```
