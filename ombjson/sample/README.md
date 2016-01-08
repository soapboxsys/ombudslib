# Sample Json Files

This directory contains sample routes and responses for queries. This is a WIP.
Do not rely on this data over the long term! Also note that none of these
objects actually exist in the block chain....!

## Individual Requests

A GET request to `https://[hostname]/api/endo/[txid]` yeilds the json in `endo.json`

A GET request to `https://[hostname]/api/bltn/[txid]` yeilds the json in
`bltn.json`

A GET request to `https://[hostname]/api/block/[hash]` yeilds the json in
`block.json` (This block is the peg block. It has no bltns or endos.)

## Paginated Requests

A GET request to `https://[hostname]/api/tag/sample` returns the json in
`tag.json`.

A GET request to `https://[hostname]/api/author/1KAVWKr2R...` 
returns the json in `author.json`. Note that the `endorsements` property is 
omitted if the author has not created any.

A GET request to `https://[]/api/loc/` 
