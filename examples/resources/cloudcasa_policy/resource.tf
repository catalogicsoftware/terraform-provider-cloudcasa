# Define a policy with a cron schedule
resource "cloudcasa_policy" "example_policy" {
  name = "cloudcasa_policy_example"
  timezone = "America/New_York"
  schedules = [
    {
      retention = 7,
      cron_spec = "30 0 * * MON,FRI",
      locked = false,
    }
  ]
}
