.PHONEY: dev backend frontend

dev: 
	make -C backend dev & make -C frontend dev

backend: 
	make -C backend dev
frontend: 
	make -C frontend dev