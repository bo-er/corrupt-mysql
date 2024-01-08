# corrupt-mysql
Corrupting MySQL to challenge the monitoring system.


# Usages

## Create a deadlock

```bash
./corrupt-mysql deadlock -H10.186.62.63 -P25690 -uuniverse_udb -p123
```

## Create a slow query log

```bash
./corrupt-mysql slowlog -H10.186.62.63 -P25690 -uuniverse_udb -p123
```
