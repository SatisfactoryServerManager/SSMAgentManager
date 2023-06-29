#/bin/bash

APPVersion="0.0.2"

rm -r ./release >/dev/null 2>&1
mkdir -p ./release/linux ./release/windows >/dev/null 2>&1

fyne-cross windows -app-id "com.ssm.ssmagentmanager" -app-version "${APPVersion}"
cp ./fyne-cross/bin/windows-amd64/* ./release/windows/.
fyne-cross linux -app-id "com.ssm.ssmagentmanager" -app-version "${APPVersion}"
cp ./fyne-cross/bin/linux-amd64/* ./release/linux/.

tar zcf "./SSMAgentManager-v${APPVersion}.tar.gz" release/*