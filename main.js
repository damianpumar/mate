const iterations = 1000;

const debounce = async (delay) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve()
    }, delay)
  })
}

const main = async () => {

  const { exec } = require('child_process')
  const pid = process.pid

  for (let i = 0; i < iterations; i++) {
    exec(`curl --location 'http://localhost:3000' \ --form 'title="from ${pid}"' \ --form 'director="me"'`)
    console.log(`[${pid}] Request ${i + 1} sent`)
  }
}

main()