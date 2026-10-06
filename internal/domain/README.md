# Domain

The rules. Everything that would still be true if we deleted the database and the API. 

Each package owns one concept: its models, its errors, its repository interface and the service 
that enforces its invariants, with persistence living in `stores/`, written against the 
interface the domain declares. 
