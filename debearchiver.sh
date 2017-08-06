LANG_OLD=$LANG
LANG="tr_TR.utf-8"
filename=$(LANG="tr_TR.utf-8" date --date="tomorrow" "+%Y-%m-%d-%d-%B-%Y-eksisozluk-debe.md")
echo "---" > $filename
echo "layout: post" >> $filename
echo -n "title: " >> $filename
date --date="tomorrow" "+%d %B %Y Ekşi Sözlük Debe" >> $filename
echo "data:" >> $filename
LANG=$LANG_OLD
for i in {1..$(python -c "from math import ceil; print int(ceil($(curl --retry 5 --retry-delay 5 -H "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:36.0) Gecko/20100101 Firefox/36.0" https://eksisozluk.com/basliklar/gundem 2>/dev/null|grep topic-list-description|sed 's/.*>\([0-9]*\) .*/\1/')/50.0))")}; do curl --retry 5 --retry-delay 5 -H "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:36.0) Gecko/20100101 Firefox/36.0" https://eksisozluk.com/basliklar/gundem\?p\=$i 2>/dev/null| sed -n '/<ul class="topic-list">/,/<\/ul>/p'|egrep 'href='|sed 's/.*"\([^"]*\)".*<small>\([0-9]*\).*/\2 \1/'; sleep 5; done | sort -nr -k1,1| awk '{print $2}'|sed -e 's@.*@https://eksisozluk.com&@' -e 's/a=popular/a=dailynice/' |xargs -I{} sh -c "curl --retry 5 --retry-delay 5 -H \"User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:36.0) Gecko/20100101 Firefox/36.0\" {} 2>/dev/null|sed -n -e '/<h1 id=\"title\"/p' -e '/<li data-id/,/<\/li/p'|sed -n '/h1 id/,/<footer>/p'|sed -e 's/<h1 .*data-title=\"\([^\"]*\)\".*/- entry_name: |\n    \1/' -e 's/.*<\/div><li.*data-id=\"\([^\"]*\)\".*data-author=\"\([^\"]*\)\".*data-favorite-count=\"\([^\"]*\)\".*/  entry_id:  \1\n  entry_writer: \2\n  entry_fav_count: \3/' -e 's/.*<div class.*/    entry_content: |/' -e '/<\/div>/d' -e '/<footer>/d'|awk 'BEGIN { RS = \"\"; FS = \"\n\" } { for (i = 1; i <= NF; i++) printf(\"%s%s\", \$i, i < NF ? \"543f562189f1e82beb9c177f89f67822\" : \"\n\") }';sleep 5"|sed 's/\(.*\)\(entry_fav_count: \)\([0-9]*\)\(543f562189f1e82beb9c177f89f67822\)  \(.*\)/\3 \1\2\3\4\5/'|sort -nr -k1,1|head -n 50|sed -e 's/^[0-9]* //' -e 's/543f562189f1e82beb9c177f89f67822/\n/g' >> $filename
echo "---" >> $filename
