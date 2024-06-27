# elastic-beat-hint test

Play with https://www.elastic.co/guide/en/beats/filebeat/current/configuration-autodiscover-hints.html, this CLI will help you to find right Kubernetes annotation to configure beat properties!

## Usage

Prefix `co.elastic.logs` is added.

```sh
elastic-beat-hint-test -a k1=v1 -a k2=v2 ....
elastic-beat-hint-test \
    --annotation enabled=true \
    -a exclude_lines='^{"log.level":"debug"(.*)$' \
    -a "processors.drop_event={\"when\":{\"or\":[{\"equals\": {\"log.level\": \"info\"}}]}}" \
    -a "json.message_key=message"
```

```yaml
enabled: true
excludeLines:
    - '''^{log.level:debug(.*)$'''
includeLines: []
json:
    message_key: ' message'
processors:
    - drop_event:
        when:
            or:
                - equals:
                    log.level: info
```
