#!/bin/bash

set -x

COOKIE='Cookie: mw_installer_session=b8e17ba511cfcfcee3ad72f869f7bdda; my_wiki_session=jfgql53goh48qhridkrmv9jm7sq4ibc1; my_wikiUserID=1; my_wikiUserName=Admin'

TOKEN=$(curl 'http://localhost:8080/api.php?action=query&meta=tokens&format=json' \
    -H "${COOKIE}" |
    jq -r '.query.tokens.csrftoken')

TEST_CONTENT=$(cat <<'EOF'
Muhterem dostum, Üstadım,

Bugün Birlik’te1 hayret ve teessürle (üzüntüyle) istifanızı okudum.2 Bu niçin böyle? Daha bilginin başlangıcında idik. İlerde himmetinizle (desteklerinizle) mecmuamız3 mutlaka daha güzelleşecekti. Acaba kim kusur etti? Neden gücendiniz? Eğer bu işte en ufak -ki bilmeyerek bir kusur olmuşsa- kendimi affetmeyeceğimden emin olabilirsiniz. Fakat kendimi zorluyorum. Hiçbir şey hatırıma gelmiyor. Hürmet ve saygı bir kusur sayılmaz sanıyorum. Ben veya bir başkası olsun sizi gücendirsek bile, sizin cemiyetten (Birlik dergisi topluluğu) öfkelenmenizi icap ettirmez. Ben şahsen rica ediyorum. Şahsiyetinizle birlikte teşekkül eden (oluşan) bu irfan havasını bozmayınız. Bu ay içinde Fatih Sayısı, Fatih Konferansları4 himmetinizden en güzel yemişlerini beklediği bugünlerde bizden yüz çevirmeniz bizim için bir gönül yarası oluyor. Bu topluluk hepimizindir. Hele mecmua bütün bütün sizindir. Onu değerli feyzinizden mahrum etmemenizi, sizi çok seven ve hürmet eden bir dostunuz olarak sizden rica ediyorum. Bugünlerde sıhhatim de pek bozuktur. Böyle zamanlarda insanın ne kadar kötümser olacağını tahmin buyurursunuz. Eğer şu işleri bırakmak lazım geliyorsa hep beraber bırakalım. Devamında fayda varsa yine hep beraber devam edelim olmaz mı?

Hürmet ve sevgiler.

Ali Nüzhet

Nejat Göyünç 

== Second Heading ==

This is some additional content!
EOF
)

curl -X POST "http://localhost:8080/api.php?action=edit&format=json&title=Test&summary=test%20summary" \
    -H "application/x-www-form-urlencoded" \
    -H "${COOKIE}" \
    --data-urlencode "text=${TEST_CONTENT}" \
    --data-urlencode "token=${TOKEN}"
