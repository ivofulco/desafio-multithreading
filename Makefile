.PHONY: desafio

desafio:
	@go run main.go ;

desafio-3:
	@sleep 1s ;
	@echo "Executando programa com delay" ;
	@sleep 1s ;
	@go run main.go ;
	@sleep 2s ;
	@go run main.go ;

desafio-parallel:
	@ make -j 2 desafio-3 desafio ;