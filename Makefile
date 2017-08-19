all: airports.go airlines.go

%.dat:
	wget https://raw.githubusercontent.com/jpatokal/openflights/master/data/$@

%.csv: %.dat
	cat $< | sed 's/\\N//g' | sed 's/\\\\//g' | sed 's/\\"/""/g' > $@

%.go: %.schema.yaml %.csv
	databundler -pkg openflights -schema $*.schema.yaml -data $*.csv -output $@
	gofmt -w $@

deps:
	go get github.com/mmcloughlin/databundler

clean:
	$(RM) *.dat *.csv

.PHONY: all clean deps

.PRECIOUS: %.dat
