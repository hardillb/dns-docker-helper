const Docker = require("dockerode")
const dns = require('dns').promises

const docker = new Docker({ socketPath: "/var/run/docker.sock" });

const re = /VIRTUAL_HOST=(.*)/

let cnames = [];

(async () => {

  await dns.setServers([
    '192.168.1.64',
    // process.env.DNS_SERVER
  ])

  const list = await docker.listContainers()

  for (var i = 0; i < list.length; i++) {
    const container = docker.getContainer(list[i].Id)
    const c = await container.inspect()
    for (const env of c.Config.Env) {
      if (re.test(env)) {
        const m = env.match(re)
        cnames.push(m[1])
      }
    }
  }

  const events = await docker.getEvents() //{ filters: { type: ['container'], event: ['create', 'start', 'stop', 'destroy', 'kill'] } })
  events.setEncoding('utf8');
  events.on('data', async ev => {
    const removed = []
    var eventJSON = JSON.parse(ev)
    if (!['start', 'stop'].includes(eventJSON.status)) {
      return;
    }
    console.log(eventJSON.status)
    if (eventJSON.status === 'start') {
      const container = docker.getContainer(eventJSON.id)
      const c = await container.inspect()
      for (const env of c.Config.Env) {
        if (re.test(env)) {
          const m = env.match(re)
          if (cnames.indexOf(m[1]) === -1) {
            cnames.push(m[1])
          }
        }
      }
    } else if (eventJSON.status === 'stop') {
      const container = docker.getContainer(eventJSON.id)
      const c = await container.inspect()
      for (const env of c.Config.Env) {
        if (re.test(env)) {
          const m = env.match(re)
          if (cnames.indexOf(m[1]) !== -1) {
            removed.push(cnames.splice(cnames.indexOf(m[1], 1))[0])
          }
        }
      }
    }
    console.log('active', cnames)
    console.log('removing',removed)
    for (host of cnames) {
      try {
        const exists = await dns.resolveCname(host)
      } catch (err) {
        console.log(err)
      }
    }
  })
})()