#!/bin/sh
# Deprecated

CONFIG_MODULE="config"
CONFIG_OUTPUT_FOLDER="deployment/output/configs"
PACKAGE_MODULE="package"
ZIP_OUTPUT_FOLDER="deployment/zip"
LAMBDA_MODULE="lambda"
BIN_OUTPUT_FOLDER="deployment/output/bin"
GO_CMD_FOLDER="cmd"
APIDOC_MODULE="apidoc"
API_DOC_DIST_FOLDER="deployment/api-doc/dist"
DB_MODULE="db"
MODULES=($CONFIG_MODULE $PACKAGE_MODULE $LAMBDA_MODULE $APIDOC_MODULE $DB_MODULE)

# receive target module
clear_module(){
    local MODULE=$1
    case "$MODULE" in
        "$CONFIG_MODULE")
            echo "remove output env config folder: $CONFIG_OUTPUT_FOLDER"
            rm -rf $CONFIG_OUTPUT_FOLDER
            ;;
        "$LAMBDA_MODULE")
            echo "remove output go executable folder: $BIN_OUTPUT_FOLDER"
            rm -rf $BIN_OUTPUT_FOLDER
            ;;
        "$PACKAGE_MODULE")
            echo "remove packaged zips: $ZIP_OUTPUT_FOLDER"
            rm -rf $ZIP_OUTPUT_FOLDER
            ;;
        "$APIDOC_MODULE")
            echo "remove api doc dist folder: $API_DOC_DIST_FOLDER"
            rm -rf $API_DOC_DIST_FOLDER
            ;;
        *)
            echo "not supported: $MODULE"
            ;;
    esac
}

# clear all modules or a module
clear_output(){
    local MODULE=$1
    if [ -z "$MODULE" ]; then 
        echo "clear all modules"
        clear_module "$CONFIG_MODULE"
        clear_module "$LAMBDA_MODULE"
        clear_module "$PACKAGE_MODULE"
        clear_module "$APIDOC_MODULE"
        clear_module "$DB_MODULE"
    else
        echo "clear 1 module"
        clear_module "$MODULE"
    fi
    echo ""
}

# receive target module
copy_module(){
    local MODULE=$1
    case "$MODULE" in
        "$CONFIG_MODULE")
            echo "copy env config to output folder: $CONFIG_OUTPUT_FOLDER"
            cp -r configs $CONFIG_OUTPUT_FOLDER
            ;;
        "$APIDOC_MODULE")
            echo "copy api doc dist to output folder: $API_DOC_DIST_FOLDER"
            cp -r api $API_DOC_DIST_FOLDER
            ;;
        *)
            echo "not supported: $MODULE"
            ;;
    esac
}

# copy all modules or a module
copy_output(){
    local MODULE=$1
    if [ -z "$MODULE" ]; then 
        echo "copy all modules"
        copy_module "$CONFIG_MODULE"
        copy_module "$LAMBDA_MODULE"
        copy_module "$PACKAGE_MODULE"
        copy_module "$APIDOC_MODULE"
        copy_module "$DB_MODULE"
    else
        echo "copy 1 module"
        copy_module "$MODULE"
    fi
    echo ""
}

# build lambda(s)
build_lambda(){
    local SERVICE=$1
    if [ "$SERVICE" == "local" ]; then
        echo "not allow to build local"
        return
    fi
    for folder in $(find ${GO_CMD_FOLDER} -mindepth 1 -maxdepth 1 -type d ); do
        NAME=$(basename ${folder})
        if [ -z "$SERVICE" -a "$NAME" == "local" ]; then
            echo "skip local"
        elif [ -z "$SERVICE" -o "$SERVICE" == "$NAME" ]; then
            echo "build lambda service $NAME"
            set GOARCH=amd64&& set GOOS=linux&& go build -o ${BIN_OUTPUT_FOLDER}/${NAME}/main ${GO_CMD_FOLDER}/${NAME}/main.go
            echo "--------> ${BIN_OUTPUT_FOLDER}/${NAME}/main"
        fi
    done
}

# receive target module
build_module(){
    local MODULE=$1
    local SUB_MODULE=$2
    case "$MODULE" in
        "$LAMBDA_MODULE")
            echo "build $LAMBDA_MODULE; lambda service:[ $SUB_MODULE ]"
            echo "pre build"
            clear_output "$CONFIG_MODULE"
            clear_output "$PACKAGE_MODULE"
            copy_output "$CONFIG_MODULE"
            echo "start build"
            build_lambda "$SUB_MODULE"
            ;;
        *)
            echo "not supported: $MODULE"
            ;;
    esac
}

# build all modules or a module
build_output(){
    local MODULE=$1
    local SUB_MODULE=$2
    if [ -z "$MODULE" ]; then
        echo "build all modules"
        for i in "${MODULES[@]}"
        do
            build_module "$i" "$SUB_MODULE"
        done
    else
        echo "build 1 module"
        build_module "$MODULE" "$SUB_MODULE"
    fi
    echo ""
}

# package lambda(s)
package_lambda(){
    local SERVICE=$1
    for folder in $(find ${BIN_OUTPUT_FOLDER} -mindepth 1 -maxdepth 1 -type d ); do
        NAME=$(basename ${folder})
        if [ -z "$SERVICE" -o "$SERVICE" == "$NAME" ]; then
            echo "package lambda service $NAME"
            cross-zip ${BIN_OUTPUT_FOLDER}/${NAME} ${ZIP_OUTPUT_FOLDER}/${NAME}.zip
            echo "--------> ${ZIP_OUTPUT_FOLDER}/${NAME}.zip"
        fi
    done
}

# receive target module
package_module(){
    local MODULE=$1
    local SUB_MODULE=$2
    case "$MODULE" in
        "$LAMBDA_MODULE")
            echo "package $LAMBDA_MODULE; lambda service:[ $SUB_MODULE ]"
            echo "pre package"
            clear_output "$PACKAGE_MODULE"
            mkdir $ZIP_OUTPUT_FOLDER
            echo "start package"
            package_lambda "$SUB_MODULE"
            ;;
        *)
            echo "not supported: $MODULE"
            ;;
    esac
}

# package all modules or a module
package_output(){
    local MODULE=$1
    local SUB_MODULE=$2
    if [ -z "$MODULE" ]; then
        echo "package all modules"
        for i in "${MODULES[@]}"
        do
            package_module "$i" "$SUB_MODULE"
        done
    else
        echo "package 1 module"
        package_module "$MODULE" "$SUB_MODULE"
    fi
    echo ""
}

# clear, copy, build, ...
output(){
    local WORK=$1
    local MODULE=$2
    local SUB_MODULE=$3
    case "$WORK" in
        "clear")
            echo "run $WORK"
            echo ""
            clear_output "$MODULE"
            ;;
        "copy")
            echo "run $WORK"
            echo ""
            copy_output "$MODULE"
            ;;
        "build")
            echo "run $WORK"
            echo ""
            build_output "$MODULE" "$SUB_MODULE"
            ;;
        "package")
            echo "run $WORK"
            echo ""
            package_output "$MODULE" "$SUB_MODULE"
            ;;
        *)
            echo "not supported: $WORK"
            ;;
    esac
}

# print modules and their folders
print_variables(){
    echo "module: $CONFIG_MODULE ; config folder: $CONFIG_OUTPUT_FOLDER"
    echo "module: $PACKAGE_MODULE ; zip folder: $ZIP_OUTPUT_FOLDER"
    echo "module: $LAMBDA_MODULE ; cmd folder: $GO_CMD_FOLDER ; bin folder: $BIN_OUTPUT_FOLDER"
    echo "module: $APIDOC_MODULE ; dist folder: $API_DOC_DIST_FOLDER"
    echo "module: $DB_MODULE"
}

# execute
echo ">>>> variables"
print_variables

echo ">>>> output"
output "$@"