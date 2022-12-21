#!/bin/sh

FOLDER=$PWD/deployment/output/configs
RSA=$FOLDER/jwt.rsa.pem
RSAPUB=$FOLDER/jwt.rsa.pem.pub
RSAPUBPEM=$FOLDER/jwt.rsa.pub.pem
YAML=$FOLDER/jwt.yml
TAB="   "

echo ">> prepare folder & clear data"
mkdir $FOLDER
rm $RSA
rm $RSAPUB
rm $RSAPUBPEM
rm $YAML
echo "<< done"

echo ">> generate rsa key pairs"
ssh-keygen -t rsa -m pem -f $RSA -q -N ""
RSAPUB_STR=$(ssh-keygen -e -f $RSA -O $RSAPUBPEM -m pkcs8)
echo "$RSAPUB_STR" >> $RSAPUBPEM
echo "files:"
echo "- $RSA"
echo "- $RSAPUB"
echo "- $RSAPUBPEM"
echo "<< done"

# write yaml
echo ">> write in yml"
echo "PRIVATE_KEY: |-" >> $YAML
while read -r line; do
    echo "${TAB}${line}" >> $YAML
done < $RSA
echo "PUBLIC_KEY: |-" >> $YAML
while read -r line; do
    echo "${TAB}${line}" >> $YAML
done < $RSAPUBPEM
echo "files:"
echo "- $YAML"
echo "<< done"
