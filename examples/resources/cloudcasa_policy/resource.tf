# Example policy
resource "cloudcasa_policy" "example" {
  name = "test_terraform_policy"
  timezone = "America/New_York"
  schedules = [
    {
      retention = 12,
      cron_spec = "30 0 * * MON,FRI",
      locked = false,
    }
  ]
}
