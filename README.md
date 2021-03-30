# carts-reseeder

Improves the reproducibility of experiments by reseeding carts-db with a
dumpfile.

carts-db uses a MongoDB instance at `carts-db:27017`. Data is saved to the
`cart` collection under the `data` database. The main contributor to load with
carts comes from users browsing to Sock Shop, which immediately persists an
empty cart for the associated user session.

For example, during a mixed load test (buying, browsing and news) of 245 users
with 165326 requests over five minutes, 23644 documents are created.
