SRC  = $(wildcard *.go *.xml Makefile)
APK = via.apk

$(APK): $(SRC)
	gomobile install

run: $(APK)
	adb shell monkey -p org.golang.todo.via -c android.intent.category.LAUNCHER 1
	adb logcat -c
	adb logcat GoLog:I *:S
