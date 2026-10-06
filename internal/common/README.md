# Common

Generic helpers that happen to be ours. `httputil` is a JSON HTTP client; it would work unchanged
in any other project, which is the whole test. 

Nothing here knows what the domain package is, or what the business is trying to solve. Keep it 
small: `common` is where unrelated code goes to pile up if nobody's watching, avoid `util`.
