package controllers

import (
  "fmt"
  "github.com/MongoHQ/api"
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
      fmt.Println(deployment.CurrentPrimary + " :: " + deployment.Id)
    }
  }
}

func Deployment(deploymentId string) {
  deployment, err := api.GetDeployment(deploymentId, OauthToken)

  if err != nil {
    fmt.Println("Error retrieving deployments: " + err.Error())
  } else {
    fmt.Println("=== " + deployment.Id)
    fmt.Println("  current primary:     " + deployment.CurrentPrimary)
    fmt.Println("  members:             " + strings.Join(deployment.Members, ","))
    fmt.Println("  version:             " + deployment.Version)

    if deployment.AllowMultipleDatabases { 
      fmt.Println("  multiple databases?: true")
    }
  }
}

func MongoStat(deployment_id, database_name string) {
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

    headerFormat  := "%" + strconv.Itoa(hostLength) + "s" + "%7s%7s%7s%7s%8s%8s%8s%7s%7s%7s%7s%" + strconv.Itoa(lockLength) + "s%11s%6s|%-3s%6s|%-3s%7s%7s%6s%11s\n"
    sprintfFormat := "%" + strconv.Itoa(hostLength) + "s" + "%7s%7s%7s%7s%8s%8s%8d%7d%7d%7d%7d%" + strconv.Itoa(lockLength) + "s%11d%6d|%-3d%6d|%-3d%7.0f%7.0f%6d%11s\n"

    if loopCount % 5 == 0 {
      fmt.Printf(headerFormat, "host", "insert", "query", "update", "delete", "getmore", "command", "flush", "mapped", "vsize", "res", "faults", "locked %", "idx miss %", "qr", "qw", "ar", "aw", "netIn", "netOut", "conn", "time")
    }

    now := time.Now()

    for position, mapMongoStat := range mongoStats {
      for host, stat := range mapMongoStat {
        var netIn, netOut float64
        if len(priorStat) > 0 && len(priorStat[position]) > 0 {
          netIn  = stat.NetIn - priorStat[position][host].NetIn
          netOut = stat.NetOut - priorStat[position][host].NetOut
        }

        fmt.Printf(sprintfFormat, hostRegex.ReplaceAllLiteralString(host, ""), stat.Inserts, stat.Query, stat.Update, stat.Delete, stat.Getmore, stat.Command, stat.Flushes, stat.Mapped, stat.Vsize, stat.Res, stat.Faults, stat.Locked, stat.IdxMiss, stat.Qr, stat.Qw, stat.Ar, stat.Aw, netIn, netOut, stat.Conn, now.Format("15:04:05"))
      }
    }

    priorStat = mongoStats
    loopCount += 1
  }

  err := api.DeploymentMongostat(deployment_id, database_name, OauthToken, outputFormatter)

  if err != nil {
    fmt.Println("Error: " + err.Error())
    os.Exit(1)
  }
}
