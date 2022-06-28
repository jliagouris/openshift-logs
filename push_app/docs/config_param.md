# Global Parameters:
| Param         |  Type  |                               Description                                |
|---------------|:------:|:------------------------------------------------------------------------:|
| parallelism   |  int   |  Maximum number of parallelism allowed in the application, 1 by default  |
| chan_buf_size |  int   |                Buffer size of the channels, 0 by default                 |
| push          |  int   | If the application supports push-based workflow, 1 for true, 0 for false |
| pull          |  int   | If the application supports pull-based workflow, 1 for true, 0 for false |
| push_period   | string |     Period for push operation. Format: Integer + unit (d/s/ms, etc)      |