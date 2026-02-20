.PHONY: nlp-server
nlp-server:
		docker run -p 9000:9000 --name corenlp --rm nlpbox/corenlp:latest
