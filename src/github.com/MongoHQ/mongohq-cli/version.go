package mongohq_cli

import (
  "os"
)

func Version() string {
  var suffix string
  if os.Getenv("CIRCLE_BUILD_NUM") != "" {
    suffix = os.Getenv("CIRCLE_BUILD_NUM")
  } else {
    suffix = "snowflake"
  }

  return "0.0.1-" + suffix
}
