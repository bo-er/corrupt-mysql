# corrupt-mysql
Corrupting MySQL and testing how the monitoring system reacts.


# Usages

## Create a deadlock

```bash
./corrupt-mysql deadlock -H10.186.62.63 -P25690 -uuniverse_udb -p123
```

## Create a slow query log

```bash
./corrupt-mysql slowlog -H10.186.62.63 -P25690 -uuniverse_udb -p123
```

## Create big transactions

```bash
./corrupt-mysql bt -H10.186.62.63 -P25690 -uuniverse_udb -p123 100MB
```

![example](./pics/bt.jpg)