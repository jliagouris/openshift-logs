# Global Parameters:
| Param         |  Type  |                               Description                                |
|---------------|:------:|:------------------------------------------------------------------------:|
| parallelism   |  int   |  Maximum number of parallelism allowed in the application, 1 by default  |
| chan_buf_size |  int   |                Buffer size of the channels, 0 by default                 |
| push          |  int   | If the application supports push-based workflow, 1 for true, 0 for false |
| pull          |  int   | If the application supports pull-based workflow, 1 for true, 0 for false |
| push_period   | string |     Period for push operation. Format: Integer + unit (d/s/ms, etc)      |

# Kafka Cluster Config: 
The parameters of cluster configuration are subject to specific server setup and security requirements. 
The following are what's needed to contact Kafka Server in Confluent Cloud.

| Param              |    Type    |               Description                |
|--------------------|:----------:|:----------------------------------------:|
| bootstrap.servers  |   string   |          Server address & port           |
| security.protocol  |   string   | Name of identity authentication protocol |
| sasl.mechanisms    |   string   |           Encryption algorithm           |
| sasl.username      |   string   |           Authorized user name           |
| sasl.password      |   string   |              User password               |
| session.timeout.ms | string/int |     Session timeout in milliseconds      |