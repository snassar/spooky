# Ansible Logging Overview

This document provides a comprehensive overview of logging in Ansible when running playbooks, covering basic output, verbosity levels, configuration options, and advanced logging techniques.

## Table of Contents

1. [Default Ansible Output](#default-ansible-output)
2. [Verbosity Levels](#verbosity-levels)
3. [Log File Configuration](#log-file-configuration)
4. [Structured Logging with Callback Plugins](#structured-logging-with-callback-plugins)
5. [Custom Logging with Callback Plugins](#custom-logging-with-callback-plugins)
6. [Logging Configuration Options](#logging-configuration-options)
7. [Task-Level Logging](#task-level-logging)
8. [Error Logging and Handling](#error-logging-and-handling)
9. [Performance Logging](#performance-logging)
10. [Integration with External Logging Systems](#integration-with-external-logging-systems)
11. [Best Practices](#best-practices)
12. [Environment-Specific Logging](#environment-specific-logging)

## Default Ansible Output

When you run an Ansible playbook, you get several types of output by default:

```bash
ansible-playbook playbook.yml
```

**Default output includes:**
- Task names and status (OK, CHANGED, FAILED, SKIPPED)
- Host information
- Task execution results
- Summary statistics

**Example output:**
```
PLAY [webservers] ***************************************************************

TASK [Gathering Facts] *********************************************************
ok: [web1.example.com]
ok: [web2.example.com]

TASK [Install nginx] ***********************************************************
changed: [web1.example.com]
changed: [web2.example.com]

TASK [Start nginx service] *****************************************************
ok: [web1.example.com]
ok: [web2.example.com]

PLAY RECAP *********************************************************************
web1.example.com              : ok=3    changed=1    unreachable=0    failed=0
web2.example.com              : ok=3    changed=1    unreachable=0    failed=0
```

## Verbosity Levels

Ansible provides multiple verbosity levels to control logging detail:

```bash
# Basic verbosity
ansible-playbook playbook.yml -v

# More verbose
ansible-playbook playbook.yml -vv

# Even more verbose
ansible-playbook playbook.yml -vvv

# Maximum verbosity
ansible-playbook playbook.yml -vvvv
```

**Verbosity levels:**
- `-v`: Basic task results
- `-vv`: Task results + task parameters
- `-vvv`: Task results + task parameters + connection info
- `-vvvv`: Maximum verbosity including internal Ansible operations

**Example verbose output (-vv):**
```
TASK [Install nginx] ***********************************************************
task path: /path/to/playbook.yml:5
<web1.example.com> ESTABLISH SSH CONNECTION FOR USER: ansible
<web1.example.com> SSH: EXEC ssh -C -o ControlMaster=auto -o ControlPersist=60s -o User=ansible -o ConnectTimeout=10 -o ControlPath=/home/user/.ansible/cp/ansible-ssh-%h-%p-%r web1.example.com '/bin/sh -c '"'"'echo ~ && sleep 0'"'"''
<web1.example.com> (0, b'/home/ansible\n', b'')
changed: [web1.example.com] => {
    "changed": true,
    "msg": "All packages installed successfully"
}
```

## Log File Configuration

You can configure Ansible to log to files by setting environment variables or in `ansible.cfg`:

### Configuration File Method

```ini
# ansible.cfg
[defaults]
log_path = /var/log/ansible.log
```

### Environment Variables Method

```bash
export ANSIBLE_LOG_PATH=/var/log/ansible.log
```

### Log File Content Example

```
2024-01-15 10:30:15,123 INFO     ansible-playbook 2.15.0
2024-01-15 10:30:15,124 INFO     config file = /etc/ansible/ansible.cfg
2024-01-15 10:30:15,125 INFO     configured module search path = ['/home/user/.ansible/plugins/modules', '/usr/share/ansible/plugins/modules']
2024-01-15 10:30:15,126 INFO     ansible python module location = /usr/lib/python3.9/site-packages/ansible
2024-01-15 10:30:15,127 INFO     executable location = /usr/bin/ansible-playbook
2024-01-15 10:30:15,128 INFO     python version = 3.9.0
2024-01-15 10:30:15,129 INFO     Using /etc/ansible/ansible.cfg as config file
```

## Structured Logging with Callback Plugins

Ansible supports callback plugins for custom logging formats:

### JSON Output

```bash
ansible-playbook playbook.yml --callback=json
```

**Available callback plugins:**
- `json`: JSON formatted output
- `profile_tasks`: Task timing information
- `timer`: Execution time tracking
- `log_plays`: Detailed play logging
- `selective`: Selective output based on task results

### JSON Output Example

```json
{
  "plays": [
    {
      "play": {
        "name": "webservers",
        "id": "webservers",
        "duration": {
          "start": "2024-01-15T10:30:15.123456Z",
          "end": "2024-01-15T10:30:45.654321Z"
        }
      },
      "tasks": [
        {
          "task": {
            "name": "Install nginx",
            "id": "install_nginx",
            "duration": {
              "start": "2024-01-15T10:30:20.123456Z",
              "end": "2024-01-15T10:30:25.654321Z"
            }
          },
          "hosts": {
            "web1.example.com": {
              "changed": true,
              "msg": "All packages installed successfully"
            }
          }
        }
      ]
    }
  ]
}
```

## Custom Logging with Callback Plugins

You can create custom callback plugins for specific logging needs:

### Basic Custom Callback Plugin

```python
# callback_plugins/custom_logger.py
from ansible.plugins.callback import CallbackBase
import json
import logging
from datetime import datetime

class CallbackModule(CallbackBase):
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'custom_logger'
    
    def __init__(self):
        super(CallbackModule, self).__init__()
        self.logger = logging.getLogger('ansible')
        self.logger.setLevel(logging.INFO)
        
        # Configure file handler
        handler = logging.FileHandler('/var/log/ansible/custom.log')
        formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
        handler.setFormatter(formatter)
        self.logger.addHandler(handler)
        
    def v2_playbook_on_start(self, playbook):
        self.logger.info(f"Starting playbook: {playbook._file_name}")
        
    def v2_playbook_on_play_start(self, play):
        self.logger.info(f"Starting play: {play.name}")
        
    def v2_playbook_on_task_start(self, task, is_conditional):
        self.logger.info(f"Starting task: {task.name}")
        
    def v2_runner_on_ok(self, result):
        self.logger.info(f"Task OK: {result.task_name} on {result.host}")
        
    def v2_runner_on_failed(self, result, ignore_errors=False):
        self.logger.error(f"Task FAILED: {result.task_name} on {result.host}")
        if result._result.get('stderr'):
            self.logger.error(f"Error details: {result._result['stderr']}")
        
    def v2_playbook_on_stats(self, stats):
        self.logger.info("Playbook execution completed")
        for host in stats.processed.keys():
            summary = stats.summarize(host)
            self.logger.info(f"Host {host}: {summary}")
```

### Advanced Custom Callback with Structured Logging

```python
# callback_plugins/structured_logger.py
from ansible.plugins.callback import CallbackBase
import json
import logging
from datetime import datetime

class CallbackModule(CallbackBase):
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'structured_logger'
    
    def __init__(self):
        super(CallbackModule, self).__init__()
        self.logger = logging.getLogger('ansible_structured')
        self.logger.setLevel(logging.INFO)
        
        # JSON formatter
        handler = logging.FileHandler('/var/log/ansible/structured.log')
        formatter = logging.Formatter('%(message)s')
        handler.setFormatter(formatter)
        self.logger.addHandler(handler)
        
    def _log_event(self, event_type, data):
        log_entry = {
            "timestamp": datetime.utcnow().isoformat(),
            "event_type": event_type,
            "data": data
        }
        self.logger.info(json.dumps(log_entry))
        
    def v2_playbook_on_start(self, playbook):
        self._log_event("playbook_start", {
            "playbook": playbook._file_name
        })
        
    def v2_playbook_on_play_start(self, play):
        self._log_event("play_start", {
            "play_name": play.name,
            "hosts": list(play.hosts)
        })
        
    def v2_playbook_on_task_start(self, task, is_conditional):
        self._log_event("task_start", {
            "task_name": task.name,
            "is_conditional": is_conditional
        })
        
    def v2_runner_on_ok(self, result):
        self._log_event("task_ok", {
            "task_name": result.task_name,
            "host": result.host,
            "changed": result._result.get('changed', False)
        })
        
    def v2_runner_on_failed(self, result, ignore_errors=False):
        self._log_event("task_failed", {
            "task_name": result.task_name,
            "host": result.host,
            "ignore_errors": ignore_errors,
            "error": result._result.get('msg', 'Unknown error')
        })
```

## Logging Configuration Options

### ansible.cfg Configuration

```ini
[defaults]
# Log file path
log_path = /var/log/ansible.log

# Display skipped hosts
display_skipped_hosts = True

# Display ok hosts
display_ok_hosts = True

# Show custom stats
show_custom_stats = True

# Timeout settings
timeout = 30

# Retry settings
retry_files_enabled = False

# Log level for internal Ansible logging
log_level = INFO

# Display task names
display_ok_hosts = True

# Display skipped hosts
display_skipped_hosts = True

# Display failed hosts
display_failed_hosts = True

# Show custom stats
show_custom_stats = True

# Display task timing
display_task_timing = True
```

### Environment Variables

```bash
# Log file path
export ANSIBLE_LOG_PATH=/var/log/ansible.log

# Log level
export ANSIBLE_LOG_LEVEL=INFO

# Display options
export ANSIBLE_DISPLAY_SKIPPED_HOSTS=True
export ANSIBLE_DISPLAY_OK_HOSTS=True
export ANSIBLE_DISPLAY_FAILED_HOSTS=True

# Timeout
export ANSIBLE_TIMEOUT=30
```

## Task-Level Logging

You can add logging within your playbook tasks:

### Basic Task Logging

```yaml
---
- name: Example playbook with logging
  hosts: all
  tasks:
    - name: Log task start
      debug:
        msg: "Starting configuration task"
      
    - name: Configure application
      template:
        src: app.conf.j2
        dest: /etc/app/app.conf
      register: config_result
      
    - name: Log configuration result
      debug:
        msg: "Configuration {{ 'succeeded' if config_result.changed else 'was already correct' }}"
```

### Advanced Task Logging with Conditions

```yaml
---
- name: Advanced logging example
  hosts: all
  vars:
    log_level: "INFO"
  tasks:
    - name: Log playbook start
      debug:
        msg: "Starting deployment on {{ inventory_hostname }}"
      when: log_level in ['DEBUG', 'INFO']
      
    - name: Check system status
      shell: systemctl status nginx
      register: nginx_status
      changed_when: false
      
    - name: Log nginx status
      debug:
        msg: "Nginx status: {{ nginx_status.stdout_lines[0] }}"
      when: log_level in ['DEBUG', 'INFO']
      
    - name: Deploy application
      copy:
        src: app.tar.gz
        dest: /tmp/app.tar.gz
      register: deploy_result
      
    - name: Log deployment result
      debug:
        msg: "Deployment {{ 'completed' if deploy_result.changed else 'skipped' }}"
      when: log_level in ['DEBUG', 'INFO', 'WARN']
```

## Error Logging and Handling

### Basic Error Handling

```yaml
- name: Handle errors with logging
  block:
    - name: Risky operation
      command: some_command
  rescue:
    - name: Log error
      debug:
        msg: "Error occurred in risky operation"
    - name: Continue despite error
      debug:
        msg: "Continuing with playbook"
```

### Advanced Error Handling with Detailed Logging

```yaml
- name: Advanced error handling
  block:
    - name: Perform critical operation
      command: critical_command
      register: critical_result
      
    - name: Log success
      debug:
        msg: "Critical operation completed successfully"
      when: critical_result.rc == 0
      
  rescue:
    - name: Log detailed error information
      debug:
        msg: |
          Critical operation failed:
          - Return code: {{ critical_result.rc }}
          - Stdout: {{ critical_result.stdout }}
          - Stderr: {{ critical_result.stderr }}
          
    - name: Send error notification
      uri:
        url: "http://monitoring.example.com/api/alerts"
        method: POST
        body_format: json
        body: |
          {
            "host": "{{ inventory_hostname }}",
            "error": "Critical operation failed",
            "details": "{{ critical_result.stderr }}"
          }
      delegate_to: localhost
      run_once: true
      
  always:
    - name: Log operation completion
      debug:
        msg: "Operation block completed"
```

## Performance Logging

### Enable Task Profiling

```bash
# Enable task profiling
ansible-playbook playbook.yml --callback=profile_tasks

# Or use timer callback
ansible-playbook playbook.yml --callback=timer
```

### Custom Performance Logging

```yaml
---
- name: Performance logging example
  hosts: all
  tasks:
    - name: Start performance timer
      set_fact:
        start_time: "{{ ansible_date_time.epoch }}"
        
    - name: Perform time-consuming operation
      shell: sleep 10
      
    - name: Calculate and log performance
      set_fact:
        end_time: "{{ ansible_date_time.epoch }}"
        duration: "{{ (ansible_date_time.epoch | int) - (start_time | int) }}"
        
    - name: Log performance metrics
      debug:
        msg: "Operation completed in {{ duration }} seconds"
        
    - name: Send performance metrics
      uri:
        url: "http://metrics.example.com/api/performance"
        method: POST
        body_format: json
        body: |
          {
            "host": "{{ inventory_hostname }}",
            "operation": "time_consuming_task",
            "duration": {{ duration }},
            "timestamp": "{{ ansible_date_time.iso8601 }}"
          }
      delegate_to: localhost
      run_once: true
```

## Integration with External Logging Systems

### Send Logs to External Systems

```yaml
- name: Send logs to external system
  uri:
    url: "http://log-server/api/logs"
    method: POST
    body_format: json
    body: "{{ ansible_facts | to_json }}"
  delegate_to: localhost
  run_once: true
```

### Integration with ELK Stack

```yaml
---
- name: Send logs to Elasticsearch
  hosts: all
  tasks:
    - name: Collect system information
      setup:
      register: system_facts
      
    - name: Send to Elasticsearch
      uri:
        url: "http://elasticsearch:9200/ansible-logs/_doc"
        method: POST
        body_format: json
        body: |
          {
            "timestamp": "{{ ansible_date_time.iso8601 }}",
            "host": "{{ inventory_hostname }}",
            "playbook": "{{ playbook_name }}",
            "task": "{{ task_name }}",
            "facts": {{ system_facts.ansible_facts | to_json }}
          }
      delegate_to: localhost
      run_once: true
```

### Integration with Splunk

```yaml
---
- name: Send logs to Splunk
  hosts: all
  tasks:
    - name: Send event to Splunk
      uri:
        url: "http://splunk:8088/services/collector"
        method: POST
        headers:
          Authorization: "Splunk {{ splunk_token }}"
        body_format: json
        body: |
          {
            "event": {
              "host": "{{ inventory_hostname }}",
              "source": "ansible",
              "sourcetype": "ansible:playbook",
              "index": "ansible_logs",
              "fields": {
                "playbook": "{{ playbook_name }}",
                "task": "{{ task_name }}",
                "status": "{{ task_status }}"
              }
            }
          }
      delegate_to: localhost
      run_once: true
```

## Best Practices

### 1. Use Appropriate Verbosity Levels

```bash
# Development environment
ansible-playbook playbook.yml -vv

# Production environment
ansible-playbook playbook.yml --callback=json > /var/log/ansible/$(date +%Y%m%d_%H%M%S).log
```

### 2. Configure Log Rotation

```bash
# /etc/logrotate.d/ansible
/var/log/ansible/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 ansible ansible
    postrotate
        systemctl reload rsyslog > /dev/null 2>&1 || true
    endscript
}
```

### 3. Use Structured Logging for Machine Processing

```yaml
---
- name: Structured logging example
  hosts: all
  tasks:
    - name: Log structured event
      debug:
        msg: |
          {
            "event": "task_started",
            "timestamp": "{{ ansible_date_time.iso8601 }}",
            "host": "{{ inventory_hostname }}",
            "task": "{{ task_name }}",
            "playbook": "{{ playbook_name }}"
          }
```

### 4. Implement Custom Callbacks for Specific Requirements

```python
# callback_plugins/security_logger.py
from ansible.plugins.callback import CallbackBase
import json
import logging

class CallbackModule(CallbackBase):
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'security_logger'
    
    def __init__(self):
        super(CallbackModule, self).__init__()
        self.logger = logging.getLogger('ansible_security')
        self.logger.setLevel(logging.INFO)
        
        # Security-focused handler
        handler = logging.FileHandler('/var/log/ansible/security.log')
        formatter = logging.Formatter('%(asctime)s - SECURITY - %(message)s')
        handler.setFormatter(formatter)
        self.logger.addHandler(handler)
        
    def v2_runner_on_failed(self, result, ignore_errors=False):
        # Log security-relevant failures
        if 'permission' in result._result.get('msg', '').lower():
            self.logger.warning(f"Permission failure on {result.host}: {result._result['msg']}")
            
    def v2_runner_on_ok(self, result):
        # Log security-relevant changes
        if 'changed' in result._result and result._result['changed']:
            if 'file' in result.task_name.lower() or 'permission' in result.task_name.lower():
                self.logger.info(f"Security-relevant change on {result.host}: {result.task_name}")
```

### 5. Log Sensitive Information Carefully

```yaml
---
- name: Secure logging example
  hosts: all
  vars:
    sensitive_password: "{{ vault_password }}"
  tasks:
    - name: Use sensitive data without logging it
      user:
        name: appuser
        password: "{{ sensitive_password }}"
      no_log: true  # Prevents password from appearing in logs
      
    - name: Log non-sensitive information
      debug:
        msg: "User appuser created successfully"
```

### 6. Use Conditional Logging to Reduce Noise

```yaml
---
- name: Conditional logging example
  hosts: all
  vars:
    debug_mode: false
  tasks:
    - name: Debug information
      debug:
        msg: "Detailed debug information"
      when: debug_mode | bool
      
    - name: Always log important events
      debug:
        msg: "Important event occurred"
```

### 7. Monitor Log File Sizes

```yaml
---
- name: Log monitoring
  hosts: localhost
  tasks:
    - name: Check log file size
      stat:
        path: /var/log/ansible.log
      register: log_stats
      
    - name: Alert if log file is too large
      debug:
        msg: "WARNING: Ansible log file is {{ (log_stats.stat.size / 1024 / 1024) | round(2) }} MB"
      when: log_stats.stat.size > 104857600  # 100MB
```

## Environment-Specific Logging

### Development Environment

```bash
# Verbose output for debugging
ansible-playbook playbook.yml -vv

# Log to development log file
ansible-playbook playbook.yml --callback=json > /tmp/ansible-dev.log
```

### Staging Environment

```bash
# Moderate verbosity
ansible-playbook playbook.yml -v

# Log to staging log file
ansible-playbook playbook.yml --callback=json > /var/log/ansible/staging/$(date +%Y%m%d_%H%M%S).log
```

### Production Environment

```bash
# Minimal output, maximum logging
ansible-playbook playbook.yml --callback=json > /var/log/ansible/production/$(date +%Y%m%d_%H%M%S).log

# With custom callback for production
ansible-playbook playbook.yml --callback=production_logger
```

### CI/CD Pipeline

```yaml
# .gitlab-ci.yml example
deploy:
  stage: deploy
  script:
    - ansible-playbook playbook.yml --callback=json > ansible.log
    - cat ansible.log
  artifacts:
    paths:
      - ansible.log
    expire_in: 1 week
```

## Conclusion

Effective logging in Ansible is crucial for debugging, monitoring, and auditing playbook executions. By understanding the different logging options and implementing appropriate strategies for your environment, you can gain valuable insights into your automation workflows and quickly identify and resolve issues.

Key takeaways:
- Use appropriate verbosity levels for different environments
- Implement structured logging for machine processing
- Configure log rotation and monitoring
- Create custom callbacks for specific requirements
- Handle sensitive information carefully
- Integrate with external logging systems when needed

Remember to regularly review and adjust your logging strategy based on your operational needs and feedback from your team.
