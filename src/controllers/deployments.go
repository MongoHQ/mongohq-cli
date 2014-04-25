package controllers

import (
  "fmt"
  "api"
  "strings"
  "os"
  "regexp"
  "strconv"
  "time"
)

func Deployments() {
  deployments, err := api.GetDeployments(OauthToken)

  if err != nil {
    fmt.Println("Error retrieving deployments: " + err.Error())
  } else {
    fmt.Println("=== My Deployments")
    for _, deployment := range deployments {
      fmt.Println(deployment.NameOrId() + " : " + deployment.CurrentPrimary)
    }
  }
}

func Deployment(deploymentId string) {
  deployment, err := api.GetDeployment(deploymentId, OauthToken)

  if err != nil {
    fmt.Println("Error retrieving deployments: " + err.Error())
  } else {
    fmt.Println("=== " + deployment.NameOrId())
    fmt.Println("  current primary:     " + deployment.CurrentPrimary)
    fmt.Println("  members:             " + strings.Join(deployment.Members, ","))
    fmt.Println("  version:             " + deployment.Version)

    if deployment.AllowMultipleDatabases { 
      fmt.Println("  multiple databases?: true")
    }

    fmt.Println("  == Databases")
    for _, database := range deployment.Databases {
      fmt.Println("    " + database.Name)
    }
  }
}

func DeploymentRename(deploymentId, name string) {
  _, err := api.RenameDeployment(deploymentId, name, OauthToken)

  if err != nil {
    fmt.Println("Error renaming deployment: " + err.Error())
  } else {
    fmt.Println("Renamed deployment to " + name + ".  You will need to reference it by the new name.")
  }
}

func CreateDeployment(deploymentName, databaseName, region string) {
  database, err := api.CreateDeployment(deploymentName, databaseName, region, OauthToken)

  if err != nil {
    fmt.Println("Error creating deployment: " + err.Error())
  } else {
    fmt.Println("=== Building deployment " + deploymentName + " with database " + database.Name)

    pollNewDeployment(database)
  }
}

func DeploymentMongoStat(deployment_id string) {
  hostRegex  := regexp.MustCompile(".(?:mongohq|mongolayer).com")
  loopCount  := 0
  var priorStat []map[string]api.MongoStat

  outputFormatter := func(mongoStats []map[string]api.MongoStat, err error) {
    if err != nil {
      fmt.Println("Error parsing stats: " + err.Error())
      os.Exit(1)
    }

    hostLength := 0
    lockLength := 0

    // Preformatting run
    for _, mapMongoStat := range mongoStats {
      for host, stats := range mapMongoStat {
        host = hostRegex.ReplaceAllLiteralString(host, "")
        if len(host) > hostLength {
          hostLength = len(host)
        }

        if len(stats.Locked) > lockLength {
          lockLength = len(stats.Locked) + 1
        }
      }
    }

    headerFormat  := "%" + strconv.Itoa(hostLength) + "s" + "%7s%7s%7s%7s%8s%8s%8s%8s%8s%8s%7s%" + strconv.Itoa(lockLength) + "s%11s%6s|%-3s%6s|%-3s%7s%7s%6s%11s\n"
    sprintfFormat := "%" + strconv.Itoa(hostLength) + "s" + "%7s%7s%7s%7s%8s%8s%8d%8s%8s%8s%7d%" + strconv.Itoa(lockLength) + "s%11d%6d|%-3d%6d|%-3d%7s%7s%6d%11s\n"

    if loopCount % 5 == 0 {
      fmt.Printf(headerFormat, "host", "insert", "query", "update", "delete", "getmore", "command", "flush", "mapped", "vsize", "res", "faults", "locked %", "idx miss %", "qr", "qw", "ar", "aw", "netIn", "netOut", "conn", "time")
    }

    now := time.Now()

    for _, mapMongoStat := range mongoStats {
      for host, stat := range mapMongoStat {
        fmt.Printf(sprintfFormat, hostRegex.ReplaceAllLiteralString(host, ""), stat.Inserts, stat.Query, stat.Update, stat.Delete, stat.Getmore, stat.Command, stat.Flushes, stat.PrettyMapped(), stat.PrettyVsize(), stat.PrettyRes(), stat.Faults, stat.Locked, stat.IdxMiss, stat.Qr, stat.Qw, stat.Ar, stat.Aw, stat.PrettyNetIn(), stat.PrettyNetOut(), stat.Conn, now.Format("15:04:05"))
      }
    }

    priorStat = mongoStats
    loopCount += 1
  }

  err := api.DeploymentMongostat(deployment_id, OauthToken, outputFormatter)

  if err != nil {
    fmt.Println("Error: " + err.Error())
    os.Exit(1)
  }
}

func DeploymentOplog(deployment_id string) {
  outputFormatter := func(entry string, err error) {
    fmt.Println(entry)
  }

  err := api.DeploymentOplog(deployment_id, OauthToken, outputFormatter)
  if err != nil {
    fmt.Println("Error: " + err.Error())
    os.Exit(1)
  }
}

func pollNewDeployment(databaseRecord api.Database) {
  var database api.Database // Schope database just in case
  var err error
  status, database := databaseRecord.Status, databaseRecord

  for status == "new" {
    fmt.Print(".")
    database, err = api.GetDatabase(database.Name, OauthToken)
    if err != nil {
      fmt.Println(err.Error())
      fmt.Println("\nError pulling database information.  For a manual update, please run:\n\n mongohq databases:info --database " + database.Name)
      os.Exit(1)
    }
    status = database.Status
  }

  deployment, err := api.GetDeployment(database.Deployment_id, OauthToken)

  if err != nil {
    fmt.Println("\nError pulling new deployment information.  For a manual update, please run:\n\n mongohq databases:info --database " + database.Name)
    os.Exit(1)
  }

  fmt.Print("\n")
  fmt.Println("Your database is ready. To add a user to your database, run:")
  fmt.Println("  mongohq users:create --deployment " + database.Deployment_id + " --database " + database.Name + " -u <username>")
  fmt.Println("")
  fmt.Println("To connect to your database, run:")
  fmt.Println("  mongo " + deployment.CurrentPrimary + "/" + database.Name + " -u <username>" + " -p")
  fmt.Println("")
  fmt.Println("Your applications should use the following URI to connect:")
  fmt.Println("  mongodb://<username>:<password>@" + strings.Join(deployment.Members, ",") + "/" + database.Name)
  fmt.Println("\nEnjoy!")
}
