SRC     = $(wildcard *.go *.xml Makefile)
SRCJAR  = via-sources/via-sources.jar
AAR     = via-aar/via.aar

default: $(AAR) $(SRCJAR)

test:
	go test

$(AAR): via.aar
	cp via.aar $(AAR)

$(SRCJAR): via.aar
	cp via-sources.jar $(SRCJAR)

via.aar: $(SRC)
	gomobile bind

$(APK): $(SRC)
	gomobile install

run: $(APK)
	adb shell monkey -p org.golang.todo.via -c android.intent.category.LAUNCHER 1
	adb logcat -c
	adb logcat GoLog:I *:S
