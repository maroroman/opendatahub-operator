MYREPO=$1
MYBRANCH=$2
 
git submodule--helper set-url odh-manifests ${MYREPO}
git submodule--helper set-branch odh-manifests --branch=${MYBRANCH}
git submodule sync
cd odh-manifests
git fetch --all
git checkout ${MYBRANCH}
cd ..
git submodule update --remote