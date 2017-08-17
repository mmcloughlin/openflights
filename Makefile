all: airports.go

%.dat:
	wget https://raw.githubusercontent.com/jpatokal/openflights/master/data/$@

%.go: make_datafile.go %.schema.yaml %.dat
	go run $< -schema $*.schema.yaml -data $*.dat -output $@
	gofmt -w $@

clean:
	$(RM) *.dat

.PHONY: all clean
