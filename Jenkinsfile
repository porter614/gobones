def SUPPORTED_PR_TAGS = [
   "fix",
   "update",
   "new",
   "breaking",
   "docs",
   "build",
   "upgrade",
   "chore",
   "merge"  // Do not use
]

def DOWNSTREAMS = [
    "bedl",
    "papr",
    "patdist",
    "rpa",
    "pdc"
]

node('mach3-porter614-non-prod-slave') {

    // Stage variables
    def currentStage
    def currentJob

    // Set environment vars
    APP = getAppFromJobName("${JOB_NAME}")
    echo "App name ${APP}"
    CLUSTER_NAME = "sat1"
    echo "Cluster name ${CLUSTER_NAME}"

    if ( "${BRANCH_NAME}" == 'master' ) {
        BRANCH = 'master'
        REF_SPEC = '+refs/heads/master:refs/remotes/origin/master' 
    } else {
        BRANCH = "origin/pr/${CHANGE_ID}/merge"
        REF_SPEC = "+refs/pull/${CHANGE_ID}/*:refs/remotes/origin/pr/${CHANGE_ID}/*"
    }

    // Get PR tag
    checkout scm
    COMMIT_MSG = sh (
        script: "git log --oneline --no-decorate -1",
        returnStdout: true
    ).trim()
    String prTag = getTagFromCommitMsg("${COMMIT_MSG}")
    if ( SUPPORTED_PR_TAGS.contains(prTag) == false ) {
        echo "Unsupported PR Tag: " + prTag
        currentBuild.result = 'FAILURE'
    }

    try {

    if ( "${BRANCH_NAME}" == 'master' ) {
      stage('Roll Version') {
          currentStage = "Roll Version"

         // Roll version based on commit
         ROLL_COMPONENT = getVersionRollFromTag(prTag)
         if ( "${ROLL_COMPONENT}" != 'none' ) {
            currentJob = build job: '../gobones_roll_version', propagate: false, parameters: [
               string(name: "APP", value: "${APP}"),
               string(name: "ROLL_COMPONENT", value: "${ROLL_COMPONENT}")]
            if ( currentJob.result != "SUCCESS" ) { throw new Exception("${currentStage} failed") }
         }
      }
    }

    stage('Build & Test') {
        currentStage = "Build & Test"
        currentJob = build job: '../gobones_unittest', propagate: false, parameters: [
            string(name: "BRANCH", value: "${BRANCH}"),
            string(name: "REF_SPEC", value: "${REF_SPEC}"),
            string(name: "APP", value: "${APP}"),
            string(name: "WORKDIR", value: "")]
        if ( currentJob.result != "SUCCESS" ) { throw new Exception("${currentStage} failed") }
    }

    // CI for PR updates stops here
    // Remaining stages are for merges to master only

    if ( "${BRANCH_NAME}" == 'master' ) {
        stage('Promote to dev') {
            currentStage = "Promote to dev"
            currentJob = build job: '../gobones_ship_image', propagate: false, parameters: [
                string(name: "BRANCH", value: "${BRANCH}"),
                string(name: "REF_SPEC", value: "${REF_SPEC}"),
                string(name: "APP", value: "${APP}"),
                string(name: "ENV", value: "dev")]
            if ( currentJob.result != "SUCCESS" ) { throw new Exception("${currentStage} failed") }
        }

        stage('CAT') {
            currentStage = "CAT"
            currentJob = build job: '../gobones_CAT', propagate: false, parameters: [
                string(name: "APP", value: "${APP}"),
                string(name: "CLUSTER_NAME", value: "${CLUSTER_NAME}")]
            if ( currentJob.result != "SUCCESS" ) { throw new Exception("${currentStage} failed") }
        }

        stage('Promote to int') {
            currentStage = "Promote to int"
            currentJob = build job: '../gobones_ship_image', propagate: false, parameters: [
                string(name: "BRANCH", value: "${BRANCH}"),
                string(name: "REF_SPEC", value: "${REF_SPEC}"),
                string(name: "APP", value: "${APP}"),
                string(name: "ENV", value: "int")]
            if ( currentJob.result != "SUCCESS" ) { throw new Exception("${currentStage} failed") }
        }

        if ( "${APP}" == 'gobones' ) {
            stage('Update Downstreams') {
                for (DS in DOWNSTREAMS) {
                    try {
                        build job: "../gobones_update_downstream", parameters: [
                            string(name: "APP", value: "${DS}")]
                    } catch (Exception ex) {
                        // Downstream updates are best effort
                        continue
                    }
                }
            }
        }
    }

    } catch (Exception ex) {
        echo "Pipeline Failure: " + ex.getMessage() 
        currentBuild.result = 'FAILURE'
    } finally {

    if ( "${BRANCH_NAME}" == 'master' ) {
        branchDescr = "master branch"
    } else {
        branchDescr = "<https://github.com/porter614/${APP}/pull/${CHANGE_ID}|PR-${CHANGE_ID}>" + 
            " (opened by ${CHANGE_AUTHOR})"
    }
              

    }   // finally
}

// Method to get a tag from a title
String getTagFromCommitMsg(String title) {
   return title.tokenize()[1].replaceAll(":","").toLowerCase()
}

String getAppFromJobName(String jn) {
  return jn.split("/")[-2].split("_")[0].toLowerCase()
}

// Method to get version roll spec from PR tag
String getVersionRollFromTag(String tag) {
   switch (tag) {
      case "fix":
      case "update":
         return "patch"
      case "new":
         return "minor"
      case "breaking":
         return "major"
      default:
         return "none"
   }
}
