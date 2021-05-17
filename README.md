# iot-dash

# Seeded Data

**Account**
Id: `testacct-0000-0000-0000-000000000000`
Username: `admin@gmail.com`
Password: `scoobydoo`
Bearer Token: `BEcSG4MsgCOk4+UTA8gtZw==`

# Running
**Pre-requites**
* Docker / Docker Compose
* Openssl
* Make

**Run**  
`make run`

# Fake IOT Testing

**Compliance**
```bash
fakeiot --token=BEcSG4MsgCOk4+UTA8gtZw== --url=https://localhost:8000 --ca-cert=/<project_path>/server/ssl/cert.pem test
```

**Generate Logins**
```bash
fakeiot --token=BEcSG4MsgCOk4+UTA8gtZw== --url=https://localhost:8000 --ca-cert=/<project_path>/server/ssl/cert.pem run --account-id="testacct-0000-0000-0000-000000000000" --period=10s --freq=0.01s --users=100
```

**Upgrade Account**
```bash
curl -v -k -H "Authorization:Bearer BEcSG4MsgCOk4+UTA8gtZw==" -XPOST https://localhost:8000/account/upgrade
```

**Get Metrics**
```bash
curl -v -k -H "Authorization:Bearer BEcSG4MsgCOk4+UTA8gtZw==" https://localhost:8000/metrics
```
