- better error reporting
- exported types cleanup
- general code cleanup
- file names
- error reporting for failed batches
- implement grapher
- implement graph indexer
- documentation
- each UrlTemplate, for objects created locally with New[NeoNode|NeoRelationship|..] can have predefined values
  which should allow using it in batches, in form like {1}/properties;
  or actually what needs to be done is to make setters from all UrlTemplate properties of NeoNode|NeoRelationship etc.
