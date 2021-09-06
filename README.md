# searchhugo
A minimalist search indexer and javascript library for hugo static websites which is as smart as modern search engines!

- **Synonyms and Antonyms**  
The search indexer uses synonyms and antonyms to find all possible ways a word might have appeared in the contents of a hugo website! That means, if your user is searching for `strong` and your text contains `powerful`, it will still appear in search results!
- **Smart rating system**  
Position of the words found in text, their vicinity to each other, their location in title or in body, the exact word, its synonym or its antonym, all have effect on the score and the appearance of the post in search results.
- **Minimal javascript library**  
A very simple javascript library is also provided to use the generated search index to provide search for the website. As well as example html file.

Currently only english language is supported for smart search index.
Smart search index is using the WordNet thesaurus to find synonyms and antonyms. See [LICENSE-thesaurus-wordnet.txt](LICENSE-thesaurus-wordnet.txt).
