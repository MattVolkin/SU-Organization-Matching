import './App.css'

function App() {
  return (
    <div className="App">
      <div className="w-3/4 m-auto">
        <div className="mt-20">
          {houses.map((d) => (
            <div className="bg-white h-[450px] text-black rounded-xl">
              <div>
                <img src={d.img} alt=""/>
              </div>

              <div>
                <p>{d.name}</p>
                <p>{d.desc}</p>
                <button>choose this option</button>
              </div>
            </div>
          ))}
        </div>
        <h2>smaller test</h2>



      </div>



      <h1>Hello World!</h1>
    </div>

  );
}

const houses = [
  {
    name: 'Lion',
    img: '/images/lion.jpg',
    desc: 'Do you feel like a Lion'
  },
  {
    name: 'Raven',
    img: '/images/raven.jpg',
    desc: 'Do you feel like a Raven'
  },
  {
    name: 'Badger',
    img: '/images/badger.jpg',
    desc: 'Do you feel like a badger'
  },
  {
    name: 'snake',
    img: '/images/snake.jpg',
    desc: 'Do you feel like a snake'
  }
]

export default App;