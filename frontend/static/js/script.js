document.querySelector('#loginBtn').addEventListener('click', () => {
  const user = document.querySelector('#user').value
  const pass = document.querySelector('#pass').value

  if (!user || !pass) {
    alert('Fill all fields')
    return
  }

  // simulate login
  localStorage.setItem('user', user)
  loadDashboard()
})

function loadDashboard() {
  document.querySelector('#app').innerHTML = `
    <section id="dashboard">
      <h2>Welcome, ${localStorage.getItem('user')}</h2>
      <button id="logout">Logout</button>
      <div id="wallet">
        <h3>Wallet</h3>
        <p>BTC: 0.5</p>
        <p>ETH: 2.3</p>
      </div>
    </section>
  `
  document.querySelector('#logout').addEventListener('click', () => {
    localStorage.clear()
    location.reload()
  })
}
