# elastic-beat-hint test

Play with https://www.elastic.co/guide/en/beats/filebeat/current/configuration-autodiscover-hints.html, this CLI will help you to find right Kubernetes annotation to configure beat properties!

## Usage

Prefix `co.elastic.logs` is added.

```sh
CLI ANNOTATION_1 ANNOTATION_2 ....
cli "enabled:true" "exclude_lines: '^{"log.level":"debug"(.*)$'"  "processors.drop_event: {\"when\":{\"or\":[{\"equals\": {\"log.level\": \"info\"}}]}}" "json.message_key: message"
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
