# Search Engine

## Description

Queries a .tsv file for a given search term using fuzzy search. The results are ranked by their BM25 score.
The .tsv file is expected to be in the following format:

```
<name>TAB<score>TAB<synonyms>TAB<info1>TAB<info2>...
```

Where:

- `<name>`: The name of the entry.
- `<score>`: A numerical score associated with the entry.
- `<synonyms>`: A semicolon-separated list of synonyms for the entry.
- `<info1>, <info2>, ...`: Additional information about the entry (optional).

## Todos

- [x] Fuzzy search using q-gram qIndex and prefix edit distance
- [ ] BM25 ranking
- [x] Similarity search using word embeddings
- [x] Use standarized format for word embeddings
- [ ] Scrape data for tv-shows
  - Use IMDb for tv-shows data
  `https://datasets.imdbws.com/`
  - Use SPARQL and Wikidata for tv-show synonyms
  ```
  SELECT ?item ?itemLabel ?alias WHERE {
  # Replace wd:Q64 with the Wikidata Q-identifier you want
  VALUES ?item { wd:Q147235 }
  
  # Retrieve aliases (skos:altLabel)
  ?item skos:altLabel ?alias .
  FILTER(LANG(?alias) = "en")
  
  # Retrieve labels (to ensure we know what the base item is)
  SERVICE wikibase:label { bd:serviceParam wikibase:language "en" }
  }
  ```
- [ ] Make it accessable via ssh (terminal website)