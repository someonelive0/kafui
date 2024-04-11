# Kafui config file template


title = "Kafui"
license = "Copyright @ 2024"


[kafka]
    name = "localhost"
    brokers = [ "127.0.0.1:9092" ]
    # sasl mechanism should be empty or "SASL_PLAINTEXT",
    # if mechanism is "SASL_PLAINTEXT", then set user and password
    sasl_mechanism = ""
    user = ""
    password = ""

[zookeeper]
    hosts = [ "127.0.0.1:2181" ]
