# Clients
Anything that talks to somebody else. An adapter speaks the outside system's wire format on one side
and our types on the other, behind an interface the caller declares — so a breaking change upstream
lands in one package, and tests pass a fake instead of hitting the network. Persistence is the
exception to the folder rather than the idea: stores live beside the domain they serve, in
`domain/<name>/stores/`.
