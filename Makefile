all: airports.go airlines.go

%.dat:
	wget https://raw.githubusercontent.com/jpatokal/openflights/master/data/$@

%.csv: %.dat
	cat $< | sed 's/\\N//g' | sed 's/\\\\//g' | sed 's/\\"/""/g' > $@

%.go: make_datafile.go %.schema.yaml %.csv
	go run $< -schema $*.schema.yaml -data $*.csv -output $@
	gofmt -w $@

clean:
	$(RM) *.dat

.PHONY: all clean
