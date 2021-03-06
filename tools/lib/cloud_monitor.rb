class CloudMonitor
  def self.namespace
    "CityService-#{GO_ENV}"
  end

  def self.report_tweets(city, count)
    client = Aws::CloudWatch::Client.new
    client.put_metric_data({
      namespace: namespace, # required
      metric_data: [ # required
          {
            metric_name: "Tweets Last Hour", # required
            dimensions: [
              {
                name: "City", # required
                value: city, # required
              },
            ],
            timestamp: Time.now,
            value: count
          },
        ]
    })
  end

  def self.report_instagrams(city, count)
    client = Aws::CloudWatch::Client.new
    client.put_metric_data({
      namespace: namespace, # required
      metric_data: [ # required
          {
            metric_name: "Instagrams Last Hour", # required
            dimensions: [
              {
                name: "City", # required
                value: city, # required
              },
            ],
            timestamp: Time.now,
            value: count
          },
        ]
    })
  end
end
