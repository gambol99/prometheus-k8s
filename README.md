
## **Prometheus Kubernetes Watcher**
-----
Is a utility service used to listen to extract and inject config for kubernetes pods that are exporting prometheus metrics endpoints. It works by simply annotating the pod definition with the metrics tag. Taking the example below, lets assume you have some webapp with an nginx frontend, both of which are exporting their metrics (nginx via /status and collectd exporter), you simply add the metrics yaml as annotation, match up the ports and indicate which ports are export promethues endpoints. The service watches out for changes in pods, extracts the details and templates our the config for promethues to pick up on file discovery. In this example pod, promethues runs in one container, prometheus-k2s in another the files are shared with an emptyDir volume (check the repo for details) 

```YAML
- job_name: 'nodes'
  file_sd_configs:
  - names: [ '/etc/prometheus/targets.d/nodes.yml' ]
- job_name: 'pods'
  file_sd_configs:
  - names: [ '/etc/prometheus/targets.d/pods.yml' ]
```

The nodes are fairly easier to add, simply watching the **/api/v1/nodes** we can get a list of nodes. The pods however require additional information. Say for example you have a pod, a web app exporting some metrics, a nginx instance with

```YAML
    apiVersion: v1
    kind: ReplicationController
    metadata:
      name: webapp
    spec:
      replicas: 2
      selector:
        name: webapp
      template:
        metadata:
          labels:
            name: webapp
          annotations:
            metrics: |
            - name: webapp
              port: 8080
            - name: nginx
              port: 8081
        spec:
          containers:
            - name: webapp
              image: gambol99/myweb-app:v1.0
              ports:
              - containerPort: 8080
            - name: nginx
              image: nginx-with-collcted
              ports:
              - containerPort: 8081
            - name: collectd
```            

### **Example Pod**
-----------------------

```YAML
    ---
    apiVersion: v1
    kind: ReplicationController
    metadata:
      name: prometheus
    spec:
      replicas: 1
      selector:
        name: prometheus
      template:
        metadata:
          labels:
            name: prometheus
        spec:
          containers:
          - name: prometheus-k8s
            image: gambol99/prometheus-k8s
            args:
              - -config=/etc/prometheus/targets.d
              - -bearer-token-file=/etc/tokens/node-register.token
              - -api-protocol=https
              - -logtostderr=true
              - -v=3 
              - -nodes=false
            volumeMounts:
            - name: targets
              mountPath: /etc/prometheus/targets.d
            - name: tokens
              mountPath: /etc/tokens
          - name: prometheus
            image: gambol99/prometheus
            ports:
            - containerPort: 9090
            args:
              - -config.file=/etc/prometheus/prometheus.yml
              - -storage.local.path=/prometheus
              - -web.console.libraries=/etc/prometheus/console_libraries
              - -web.console.templates=/etc/prometheus/console
            volumeMounts:
            - name: targets
              mountPath: /etc/prometheus/targets.d
          imagePullPolicy: Always
          volumes:
          - name: targets
            source:
              emptyDir: {}
          - name: tokens
            hostPath:
              path: /run/kube-kubelet
```
