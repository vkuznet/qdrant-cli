# qdrant-cli
A simple toolkit to work with Qdrant database in terminal shell

```
# Build
make

# List collections
qdrant-cli --list

# Create collection
qdrant-cli --create my_collection --dim 384

# Create + insert dummy points
qdrant-cli --create test --dim 4 --insert-dummy 20

# Scroll records
qdrant-cli --scroll test --limit 5

# Delete collection
qdrant-cli --delete test
```
