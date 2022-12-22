#!/bin/sh
# -a account reuiqred
# -u username required
# -p profile
# -c code required
# -h help

AWS_VERION=$(aws --version)

if [ $? -ne 0 ]; then
  echo "AWS CLI is not installed; exiting"
  exit 1
else
  echo "Using AWS CLI VERSION $AWS_VERION"
fi

ACCOUNT=""
USERNAME=""
PROFILE=""
CODE=""
while getopts ":a:u:p:c:h" optname
do
    case "$optname" in
      "a")
        ACCOUNT=$OPTARG
        echo "get option -a(acount),  value is $ACCOUNT"
        ;;
      "u")
        USERNAME=$OPTARG
        echo "get option -u(username),value is $USERNAME"
        ;;
      "p")
        PROFILE=$OPTARG
        echo "get option -p(profile), value is $PROFILE"
        ;;
      "c")
        CODE=$OPTARG
        echo "get option -c(mfa code),value is $CODE"
        ;;
      "h")
        echo "usage: ./generator.sh [-h] [-a {aws account}] [-u {aws iam username}] [-p {profile name}] [-c {mfacode}]"
        echo "positional arguments:"
        echo "  -a      AWS Account"
        echo "  -u      IAM username"
        echo "  -p      AWS Profile used in issuing the temp token"
        echo "  -c      Multi-factor authentication (MFA) code"
        echo "  -h      Show this help message and exit"
        exit
        ;;
    esac
done

if [ -z "$ACCOUNT" ]; then
    echo "Account is required"
    exit 1
fi
if [ -z "$USERNAME" ]; then
    echo "Username is required"
    exit 1
fi
if [ -z "$PROFILE" ]; then
    echo "Profile is required"
    exit 1
fi
if [ -z "$CODE" ]; then
    echo "MFA code is required"
    exit 1
fi

AWS_ARN_MFA="arn:aws:iam::${ACCOUNT}:mfa/${USERNAME}"
SESSION_PROFILE="default"
if [[ ! -z "$PROFILE" ]]; then
    SESSION_PROFILE="${PROFILE}_default"
fi

AWS_STS_RESULT=$(aws sts get-session-token --duration-seconds 129600 \
  --serial-number $AWS_ARN_MFA \
  --profile $PROFILE \
  --token-code $CODE \
  --output text)
read AWS_RETURN_TYPE AWS_ACCESS_KEY_ID AWS_EXPIRATION AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN <<< $AWS_STS_RESULT
if [ -z "$AWS_ACCESS_KEY_ID" ]; then
  exit 1
fi

echo "---------------"
echo "ArnMFA:          $AWS_ARN_MFA" 
echo "SessionProfile:  $SESSION_PROFILE" 
echo "AccessKeyId:     $AWS_ACCESS_KEY_ID" 
echo "SecretAccessKey: $AWS_SECRET_ACCESS_KEY" 
echo "SessionToke:     $AWS_SESSION_TOKEN" 
echo "Expiration:      $AWS_EXPIRATION" 
echo "---------------"

echo "updating credentials of $SESSION_PROFILE"
aws configure set aws_access_key_id "$AWS_ACCESS_KEY_ID" --profile $SESSION_PROFILE 
aws configure set aws_secret_access_key "$AWS_SECRET_ACCESS_KEY"  --profile $SESSION_PROFILE 
aws configure set aws_session_token "$AWS_SESSION_TOKEN"  --profile $SESSION_PROFILE 
echo "done"