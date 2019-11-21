package cmd

import(
   "fmt"
   "github.com/spf13/cobra"
   "github.com/porter614/gobones/pkg/version"
)

var VersionCmd = &cobra.Command{
   Use: "version",
   Short: "Print version",
   Long: "Prints out the version and metadata of this go module",
   Run: printVersion,
}

var SemVerCmd = &cobra.Command{
   Use: "semver",
   Short: "Print semver",
   Long: "Prints out the semantic version of this go module",
   Run: printSemVer,
}

func init() {
   RootCmd.AddCommand(VersionCmd)
   RootCmd.AddCommand(SemVerCmd)
}

func printSemVer(cmd *cobra.Command, args []string) {
   fmt.Println("SemVer:", version.SemVer)
}

func printVersion(cmd *cobra.Command, args []string) {
   fmt.Println("****************************************************************")
   fmt.Println("Application:", version.App)
   fmt.Println("Version:", version.Version)
   fmt.Println("SemVer:", version.SemVer)
   fmt.Println("Build Info:")
   fmt.Printf("    GIT Branch: %s\n", version.Branch)
   fmt.Printf("    GIT Commit ID: %s\n", version.CommitId)
   fmt.Printf("    Build User: %s\n", version.BuildUser)
   fmt.Printf("    Build Date: %s\n", version.BuildDate)
   fmt.Println("****************************************************************")
}
