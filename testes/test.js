;(async () => {
  const fs = require('fs').promises

  const size = 24
  const xteaKey = [3727812935, 1206674243, 4196718677, 2568959323]

  const buffer = await fs.readFile('./buffer_antes')
  console.log(buffer)

  const encrypted = xteaEncrypt(buffer, size, xteaKey)

  console.log(encrypted)
  const bufferAfter = await fs.readFile('./buffer_depois')
  console.log('\n[!] Original:')
  console.log(bufferAfter)
  console.log('Are equal:', !Buffer.compare(bufferAfter, encrypted))
})()

function xteaEncrypt(buffer, size, key) {
  let u32 = new Uint32Array(buffer.buffer, 0, size)

  const delta = 0x9e3779b9
  const num_rounds = 32
  
  console.log(u32.length)

  for (let i = 2; i < u32.length / 4; i += 2) {
    u32[0] = 0 // sum

    for (let j = 0; j < num_rounds; j++) {
      u32[i] += (((u32[i + 1] << 4) ^ (u32[i + 1] >>> 5)) + u32[i + 1]) ^ (u32[0] + key[u32[0] & 3])

      u32[0] += delta

      u32[i + 1] += (((u32[i] << 4) ^ (u32[i] >>> 5)) + u32[i]) ^ (u32[0] + key[(u32[0] >> 11) & 3])
    }
  }

  return Buffer.from(u32.buffer)
}
