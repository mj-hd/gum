all: main

main: go_packages
	gom build main.go

go_packages:
	gom install

clean:
	rm -f main
