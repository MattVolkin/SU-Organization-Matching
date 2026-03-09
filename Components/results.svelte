  <h1>{results[pageNum-1]}</h1>
  <h2>Hi we are {results[pageNum-1]} and we are commited to to provide a safe space to play games and hang out with other computer nerds </h2>

  <button onclick={prevPage} disabled={pageNum === 1}>Previous</button>
  <button onclick={nextPage}>Next</button>
  <h3> Activities we do include</h3>
  <ul>
    <li>Playing games</li>
    <li>Trivia nights </li>
    <li>Presentation nights</li>
    <li>Video game tournoments</li>
    <li>Scavenger hunts</li>
  </ul>
  <h3>Meeting Information </h3>
  <p> Every Thursday at 6:30 pm in FJS 310 (the CS lounge)</p>