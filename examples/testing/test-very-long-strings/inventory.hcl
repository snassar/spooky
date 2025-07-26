# Inventory with very long strings

inventory {
  machine "long-string-server" {
    host     = "192.168.1.100"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
      description = "This is a very long description that contains many characters and should test the system's ability to handle long strings in configuration files. It includes various types of content like URLs, paths, and descriptive text that might be encountered in real-world scenarios. The string continues with more content to ensure we reach a significant length that could potentially cause issues with parsing, validation, or storage. We want to make sure the system can handle strings of this length without problems."
      long_url = "https://very-long-url.example.com/with/many/path/segments/and/query/parameters?param1=value1&param2=value2&param3=value3&param4=value4&param5=value5&param6=value6&param7=value7&param8=value8&param9=value9&param10=value10"
      long_path = "/very/long/path/to/some/file/that/might/exist/on/the/system/and/could/be/referenced/in/configuration/files/with/many/subdirectories/and/deeply/nested/structure"
    }
  }
}